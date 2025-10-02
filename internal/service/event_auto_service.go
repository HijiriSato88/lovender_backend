package service

import (
	"context"
	"fmt"
	"log"
	"lovender_backend/internal/client"
	"lovender_backend/internal/models"
	"lovender_backend/internal/repository"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// EventAutoService イベント自動登録サービス
type EventAutoService struct {
	eventsRepo     repository.EventsRepository
	keywordCache   *KeywordCacheService
	externalClient *client.ExternalPostClient
}

// NewEventAutoService コンストラクタ
func NewEventAutoService(
	eventsRepo repository.EventsRepository,
	keywordCache *KeywordCacheService,
) *EventAutoService {
	return &EventAutoService{
		eventsRepo:     eventsRepo,
		keywordCache:   keywordCache,
		externalClient: client.NewExternalPostClient(),
	}
}

// 全ユーザーの推しから自動イベント作成を実行
func (s *EventAutoService) ProcessAutoEventCreation(ctx context.Context) (*AutoEventResult, error) {
	log.Println("Starting auto event creation process")

	// 全推し情報を取得
	oshis, err := s.eventsRepo.GetAllOshisWithAccountsAndCategories()
	if err != nil {
		return nil, fmt.Errorf("failed to get all oshis: %w", err)
	}

	log.Printf("Found %d oshis to process", len(oshis))

	result := &AutoEventResult{
		ProcessedOshis: 0,
		CreatedEvents:  0,
		Errors:         []string{},
	}

	// 並列処理用のチャネルとワーカープール
	const maxWorkers = 10
	oshiChan := make(chan *models.OshiWithDetails, len(oshis))
	resultChan := make(chan *OshiProcessResult, len(oshis))

	// ワーカーを起動
	var wg sync.WaitGroup
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go s.processOshiWorker(ctx, oshiChan, resultChan, &wg)
	}

	// 推しをチャネルに送信
	for _, oshi := range oshis {
		select {
		case oshiChan <- oshi:
		case <-ctx.Done():
			close(oshiChan)
			return nil, ctx.Err()
		}
	}
	close(oshiChan)

	// ワーカー完了を待つ
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 結果を集計
	for oshiResult := range resultChan {
		result.ProcessedOshis++
		result.CreatedEvents += oshiResult.CreatedEvents
		if oshiResult.Error != "" {
			result.Errors = append(result.Errors, oshiResult.Error)
		}
	}

	log.Printf("Auto event creation completed. Processed: %d oshis, Created: %d events, Errors: %d",
		result.ProcessedOshis, result.CreatedEvents, len(result.Errors))

	return result, nil
}

// 推し処理ワーカー
func (s *EventAutoService) processOshiWorker(
	ctx context.Context,
	oshiChan <-chan *models.OshiWithDetails,
	resultChan chan<- *OshiProcessResult,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for oshi := range oshiChan {
		select {
		case <-ctx.Done():
			return
		default:
			result := s.processOshiPosts(ctx, oshi)
			resultChan <- result
		}
	}
}

// 推しの投稿を処理
func (s *EventAutoService) processOshiPosts(ctx context.Context, oshi *models.OshiWithDetails) *OshiProcessResult {
	result := &OshiProcessResult{
		OshiID:        oshi.Oshi.ID,
		OshiName:      oshi.Oshi.Name,
		CreatedEvents: 0,
		Error:         "",
	}

	// アカウントがない場合はスキップ
	if len(oshi.Accounts) == 0 {
		return result
	}

	// カテゴリIDを抽出
	var categoryIDs []uint16
	for _, category := range oshi.Categories {
		categoryIDs = append(categoryIDs, category.ID)
	}

	// カテゴリに関連するキーワードを取得
	keywords := s.keywordCache.GetKeywordsByCategories(categoryIDs)
	if len(keywords) == 0 {
		return result
	}

	// 各アカウントの投稿を処理
	for _, account := range oshi.Accounts {
		accountName := s.extractAccountName(account.URL)
		if accountName == "" {
			continue
		}

		// 最新5件の投稿を取得
		posts, err := s.externalClient.GetLatestPostsByUsername(accountName, 5)
		if err != nil {
			result.Error = fmt.Sprintf("Failed to get posts for %s: %v", accountName, err)
			continue
		}

		log.Printf("Fetched %d posts for account %s", len(posts), accountName)

		// 投稿を処理
		for _, post := range posts {
			select {
			case <-ctx.Done():
				return result
			default:
				if s.processPost(oshi.Oshi.ID, post, keywords) {
					result.CreatedEvents++
				}
			}
		}
	}

	return result
}

// 投稿を処理してイベント作成
func (s *EventAutoService) processPost(oshiID int64, post models.ExternalPost, keywords []repository.CategoryKeyword) bool {
	// 既に登録済みかチェック
	exists, err := s.eventsRepo.CheckEventExistsByPostID(post.ID)
	if err != nil {
		log.Printf("Failed to check event existence for post %d: %v", post.ID, err)
		return false
	}
	if exists {
		return false
	}

	// キーワードマッチング
	var matchedKeywords []string
	var matchedCategoryID *uint16

	content := strings.ToLower(post.Content)
	for _, keyword := range keywords {
		if strings.Contains(content, strings.ToLower(keyword.Keyword)) {
			matchedKeywords = append(matchedKeywords, keyword.Keyword)
			if matchedCategoryID == nil {
				matchedCategoryID = &keyword.CategoryID
			}
		}
	}

	// キーワードが一致しない場合はスキップ
	if len(matchedKeywords) == 0 {
		return false
	}

	// 投稿日時をパース
	createdAt, err := time.Parse("2006-01-02 15:04:05", post.CreatedAt)
	if err != nil {
		log.Printf("Failed to parse created_at for post %d: %v", post.ID, err)
		return false
	}

	// 投稿内容から日時情報を抽出
	startsAt, endsAt := s.extractDateTime(post.Content, createdAt)

	// イベントタイトル生成
	title := fmt.Sprintf("%sの投稿情報", post.User.Name)

	log.Printf("Post[%d] - Event time: %s to %s",
		post.ID,
		startsAt.Format("2006-01-02 15:04:05"),
		func() string {
			if endsAt != nil {
				return endsAt.Format("2006-01-02 15:04:05")
			}
			return "nil"
		}())

	// イベント作成
	err = s.eventsRepo.CreateAutoEvent(
		oshiID,
		post.ID,
		title,
		post.Content,
		matchedCategoryID,
		startsAt,
		endsAt,
	)
	if err != nil {
		log.Printf("Failed to create auto event for post %d: %v", post.ID, err)
		return false
	}

	log.Printf("Created auto event for oshi %d, post %d, keywords: %v",
		oshiID, post.ID, matchedKeywords)
	return true
}

// URLからアカウント名を抽出
func (s *EventAutoService) extractAccountName(url string) string {
	// 最後のスラッシュ以降を取得
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

// 投稿内容から日時情報を抽出
func (s *EventAutoService) extractDateTime(content string, postCreatedAt time.Time) (time.Time, *time.Time) {
	// 日時抽出用の正規表現パターン
	patterns := []struct {
		regex   *regexp.Regexp
		handler func([]string, time.Time) (time.Time, *time.Time)
	}{
		// パターン1: "2026年1月10日 14:00-16:00" (年月日+時刻範囲)
		{
			regexp.MustCompile(`(\d{4})年(\d{1,2})月(\d{1,2})日\s*(\d{1,2}):(\d{2})\s*[-〜～]\s*(\d{1,2}):(\d{2})`),
			s.handleYearDateTimeRange,
		},
		// パターン2: "2026年1月10日 14:00" (年月日+時刻)
		{
			regexp.MustCompile(`(\d{4})年(\d{1,2})月(\d{1,2})日\s*(\d{1,2}):(\d{2})`),
			s.handleYearDateTime,
		},
		// パターン3: "2026年1月10日" (年月日のみ)
		{
			regexp.MustCompile(`(\d{4})年(\d{1,2})月(\d{1,2})日`),
			s.handleYearDateOnly,
		},
		// パターン4: "10月3日 14:00-16:00" (月日+時刻範囲)
		{
			regexp.MustCompile(`(\d{1,2})月(\d{1,2})日\s*(\d{1,2}):(\d{2})\s*[-〜～]\s*(\d{1,2}):(\d{2})`),
			s.handleDateTimeRange,
		},
		// パターン5: "10月3日 14:00" (月日+時刻)
		{
			regexp.MustCompile(`(\d{1,2})月(\d{1,2})日\s*(\d{1,2}):(\d{2})`),
			s.handleDateTime,
		},
		// パターン6: "14:00-16:00" (時刻範囲のみ)
		{
			regexp.MustCompile(`(\d{1,2}):(\d{2})\s*[-〜～]\s*(\d{1,2}):(\d{2})`),
			s.handleTimeRange,
		},
		// パターン7: "14時〜16時" (時刻範囲・日本語)
		{
			regexp.MustCompile(`(\d{1,2})時\s*[〜～]\s*(\d{1,2})時`),
			s.handleJapaneseTimeRange,
		},
		// パターン8: "14時から16時" (時刻範囲・から)
		{
			regexp.MustCompile(`(\d{1,2})時から\s*(\d{1,2})時`),
			s.handleJapaneseTimeFromTo,
		},
		// パターン9: "14時から" (開始時刻のみ・から)
		{
			regexp.MustCompile(`(\d{1,2})時から[！!]?`),
			s.handleJapaneseTimeFrom,
		},
		// パターン10: "14時〜" (開始時刻のみ・〜)
		{
			regexp.MustCompile(`(\d{1,2})時[〜～][！!]?`),
			s.handleJapaneseTimeStart,
		},
		// パターン11: "14:00" (時刻のみ)
		{
			regexp.MustCompile(`(\d{1,2}):(\d{2})`),
			s.handleTimeOnly,
		},
		// パターン12: "10月3日" (月日のみ)
		{
			regexp.MustCompile(`(\d{1,2})月(\d{1,2})日`),
			s.handleDateOnly,
		},
	}

	// 各パターンを試行
	for _, pattern := range patterns {
		matches := pattern.regex.FindStringSubmatch(content)
		if len(matches) > 0 {
			log.Printf("DateTime extraction - Pattern matched: %v", matches)
			return pattern.handler(matches, postCreatedAt)
		}
	}

	// パターンが見つからない場合はデフォルト（投稿日の0:00-1:00）
	log.Printf("DateTime extraction - No pattern found, using default time")
	return s.getDefaultDateTime(postCreatedAt)
}

// 年月日+時刻範囲の処理 (例: "2026年1月10日 14:00-16:00")
func (s *EventAutoService) handleYearDateTimeRange(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	year, _ := strconv.Atoi(matches[1])
	month, _ := strconv.Atoi(matches[2])
	day, _ := strconv.Atoi(matches[3])
	startHour, _ := strconv.Atoi(matches[4])
	startMin, _ := strconv.Atoi(matches[5])
	endHour, _ := strconv.Atoi(matches[6])
	endMin, _ := strconv.Atoi(matches[7])

	startsAt := time.Date(year, time.Month(month), day, startHour, startMin, 0, 0, postCreatedAt.Location())
	endsAt := time.Date(year, time.Month(month), day, endHour, endMin, 0, 0, postCreatedAt.Location())

	return startsAt, &endsAt
}

// 年月日+時刻の処理 (例: "2026年1月10日 14:00")
func (s *EventAutoService) handleYearDateTime(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	year, _ := strconv.Atoi(matches[1])
	month, _ := strconv.Atoi(matches[2])
	day, _ := strconv.Atoi(matches[3])
	hour, _ := strconv.Atoi(matches[4])
	min, _ := strconv.Atoi(matches[5])

	startsAt := time.Date(year, time.Month(month), day, hour, min, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour) // 1時間後を終了時刻とする

	return startsAt, &endsAt
}

// 年月日のみの処理 (例: "2026年1月10日")
func (s *EventAutoService) handleYearDateOnly(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	year, _ := strconv.Atoi(matches[1])
	month, _ := strconv.Atoi(matches[2])
	day, _ := strconv.Atoi(matches[3])

	startsAt := time.Date(year, time.Month(month), day, 0, 0, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// 日付+時刻範囲の処理 (例: "10月3日 14:00-16:00")
func (s *EventAutoService) handleDateTimeRange(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	month, _ := strconv.Atoi(matches[1])
	day, _ := strconv.Atoi(matches[2])
	startHour, _ := strconv.Atoi(matches[3])
	startMin, _ := strconv.Atoi(matches[4])
	endHour, _ := strconv.Atoi(matches[5])
	endMin, _ := strconv.Atoi(matches[6])

	year := postCreatedAt.Year()
	// 月が過去の場合は翌年とする
	if month < int(postCreatedAt.Month()) {
		year++
	}

	startsAt := time.Date(year, time.Month(month), day, startHour, startMin, 0, 0, postCreatedAt.Location())
	endsAt := time.Date(year, time.Month(month), day, endHour, endMin, 0, 0, postCreatedAt.Location())

	return startsAt, &endsAt
}

// 日付+時刻の処理 (例: "10月3日 14:00")
func (s *EventAutoService) handleDateTime(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	month, _ := strconv.Atoi(matches[1])
	day, _ := strconv.Atoi(matches[2])
	hour, _ := strconv.Atoi(matches[3])
	min, _ := strconv.Atoi(matches[4])

	year := postCreatedAt.Year()
	if month < int(postCreatedAt.Month()) {
		year++
	}

	startsAt := time.Date(year, time.Month(month), day, hour, min, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour) // 1時間後を終了時刻とする

	return startsAt, &endsAt
}

// 時刻範囲の処理 (例: "14:00-16:00")
func (s *EventAutoService) handleTimeRange(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	startHour, _ := strconv.Atoi(matches[1])
	startMin, _ := strconv.Atoi(matches[2])
	endHour, _ := strconv.Atoi(matches[3])
	endMin, _ := strconv.Atoi(matches[4])

	// 投稿日の指定時刻
	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		startHour, startMin, 0, 0, postCreatedAt.Location())
	endsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		endHour, endMin, 0, 0, postCreatedAt.Location())

	return startsAt, &endsAt
}

// 日本語時刻範囲の処理 (例: "14時〜16時")
func (s *EventAutoService) handleJapaneseTimeRange(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	startHour, _ := strconv.Atoi(matches[1])
	endHour, _ := strconv.Atoi(matches[2])

	// 投稿日の指定時刻（分は0とする）
	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		startHour, 0, 0, 0, postCreatedAt.Location())
	endsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		endHour, 0, 0, 0, postCreatedAt.Location())

	return startsAt, &endsAt
}

// 日本語時刻範囲の処理 (例: "14時から16時")
func (s *EventAutoService) handleJapaneseTimeFromTo(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	startHour, _ := strconv.Atoi(matches[1])
	endHour, _ := strconv.Atoi(matches[2])

	// 投稿日の指定時刻（分は0とする）
	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		startHour, 0, 0, 0, postCreatedAt.Location())
	endsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		endHour, 0, 0, 0, postCreatedAt.Location())

	return startsAt, &endsAt
}

// 日本語開始時刻の処理 (例: "14時から")
func (s *EventAutoService) handleJapaneseTimeFrom(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	hour, _ := strconv.Atoi(matches[1])

	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		hour, 0, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour) // 1時間後を終了時刻とする

	return startsAt, &endsAt
}

// 日本語開始時刻の処理 (例: "14時〜")
func (s *EventAutoService) handleJapaneseTimeStart(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	hour, _ := strconv.Atoi(matches[1])

	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		hour, 0, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour) // 1時間後を終了時刻とする

	return startsAt, &endsAt
}

// 時刻のみの処理 (例: "14:00")
func (s *EventAutoService) handleTimeOnly(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	hour, _ := strconv.Atoi(matches[1])
	min, _ := strconv.Atoi(matches[2])

	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		hour, min, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// 日付のみの処理 (例: "10月3日")
func (s *EventAutoService) handleDateOnly(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	month, _ := strconv.Atoi(matches[1])
	day, _ := strconv.Atoi(matches[2])

	year := postCreatedAt.Year()
	if month < int(postCreatedAt.Month()) {
		year++
	}

	startsAt := time.Date(year, time.Month(month), day, 0, 0, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// デフォルト日時の取得 (投稿日の0:00-1:00)
func (s *EventAutoService) getDefaultDateTime(postCreatedAt time.Time) (time.Time, *time.Time) {
	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		0, 0, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// 自動イベント作成結果
type AutoEventResult struct {
	ProcessedOshis int      `json:"processed_oshis"`
	CreatedEvents  int      `json:"created_events"`
	Errors         []string `json:"errors,omitempty"`
}

// 推し処理結果
type OshiProcessResult struct {
	OshiID        int64  `json:"oshi_id"`
	OshiName      string `json:"oshi_name"`
	CreatedEvents int    `json:"created_events"`
	Error         string `json:"error,omitempty"`
}

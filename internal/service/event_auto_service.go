package service

import (
	"context"
	"fmt"
	"log"
	"lovender_backend/internal/client"
	"lovender_backend/internal/models"
	"lovender_backend/internal/repository"
	"strings"
	"sync"
	"time"
)

// EventAutoService イベント自動登録サービス
type EventAutoService struct {
	eventsRepo        repository.EventsRepository
	keywordCache      *KeywordCacheService
	externalClient    *client.ExternalPostClient
	dateTimeExtractor *DateTimeExtractionService
	jstLocation       *time.Location
}

// NewEventAutoService コンストラクタ
func NewEventAutoService(
	eventsRepo repository.EventsRepository,
	keywordCache *KeywordCacheService,
) *EventAutoService {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Printf("Warning: Failed to load JST location, using UTC: %v", err)
		jst = time.UTC
	}

	return &EventAutoService{
		eventsRepo:        eventsRepo,
		keywordCache:      keywordCache,
		externalClient:    client.NewExternalPostClient(),
		dateTimeExtractor: NewDateTimeExtractionService(),
		jstLocation:       jst,
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
	exists, err := s.eventsRepo.CheckEventExistsByPostIDAndOshiID(post.ID, oshiID)
	if err != nil {
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

	// 投稿日時をパース（日本時間として扱う）
	createdAt, err := time.ParseInLocation("2006-01-02 15:04:05", post.CreatedAt, s.jstLocation)
	if err != nil {
		log.Printf("Failed to parse created_at for post %d: %v", post.ID, err)
		return false
	}

	// 投稿内容から日時情報を抽出（日本時間として抽出される）
	startsAt, endsAt, hasDateTimePattern := s.dateTimeExtractor.ExtractDateTime(post.Content, createdAt)

	// 日時パターンが見つからない場合はスキップ
	if !hasDateTimePattern {
		log.Printf("Post[%d] - No datetime pattern found, skipping event creation", post.ID)
		return false
	}

	// 日時をUTCに変換
	startsAtUTC := startsAt.UTC()
	var endsAtUTC *time.Time
	if endsAt != nil {
		utcTime := endsAt.UTC()
		endsAtUTC = &utcTime
	}

	// イベントタイトル生成
	title := fmt.Sprintf(post.User.Name)

	// イベント作成（UTC時刻で保存）
	err = s.eventsRepo.CreateAutoEvent(
		oshiID,
		post.ID,
		title,
		post.Content,
		matchedCategoryID,
		startsAtUTC,
		endsAtUTC,
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

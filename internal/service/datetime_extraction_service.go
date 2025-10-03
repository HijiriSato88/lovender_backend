package service

import (
	"log"
	"regexp"
	"strconv"
	"time"
)

// DateTimeExtractionService 日時抽出サービス
type DateTimeExtractionService struct{}

// NewDateTimeExtractionService コンストラクタ
func NewDateTimeExtractionService() *DateTimeExtractionService {
	return &DateTimeExtractionService{}
}

// ExtractDateTime 投稿内容から日時情報を抽出
func (s *DateTimeExtractionService) ExtractDateTime(content string, postCreatedAt time.Time) (time.Time, *time.Time, bool) {
	// 日時抽出用の正規表現パターン
	patterns := []struct {
		regex   *regexp.Regexp
		handler func([]string, time.Time) (time.Time, *time.Time)
	}{
		// パターン1: "2026年1月10日 14:00-16:00" (年月日+時刻範囲)
		{
			regexp.MustCompile(`(\d{4})年(\d{1,2})月(\d{1,2})日[\s　]*(\d{1,2}):(\d{2})\s*[-〜～]\s*(\d{1,2}):(\d{2})`),
			s.handleYearDateTimeRange,
		},
		// パターン2: "2026年1月10日 14:00" (年月日+時刻)
		{
			regexp.MustCompile(`(\d{4})年(\d{1,2})月(\d{1,2})日[\s　]*(\d{1,2}):(\d{2})`),
			s.handleYearDateTime,
		},
		// パターン3: "2026年1月10日" (年月日のみ)
		{
			regexp.MustCompile(`(\d{4})年(\d{1,2})月(\d{1,2})日`),
			s.handleYearDateOnly,
		},
		// パターン4: "10/6（月）18:00まで" (月日曜日+時刻まで)
		{
			regexp.MustCompile(`(\d{1,2})/(\d{1,2})（[月火水木金土日]）[\s　]*(\d{1,2}):(\d{2})\s*まで`),
			s.handleSlashDateWeekdayTimeUntil,
		},
		// パターン5: "10/6（月）18:00-20:00" (月日曜日+時刻範囲)
		{
			regexp.MustCompile(`(\d{1,2})/(\d{1,2})（[月火水木金土日]）[\s　]*(\d{1,2}):(\d{2})\s*[-〜～]\s*(\d{1,2}):(\d{2})`),
			s.handleSlashDateWeekdayTimeRange,
		},
		// パターン6: "10/6（月）18:00" (月日曜日+時刻)
		{
			regexp.MustCompile(`(\d{1,2})/(\d{1,2})（[月火水木金土日]）[\s　]*(\d{1,2}):(\d{2})`),
			s.handleSlashDateWeekdayTime,
		},
		// パターン7: "10/6（月）" (月日曜日のみ)
		{
			regexp.MustCompile(`(\d{1,2})/(\d{1,2})（[月火水木金土日]）`),
			s.handleSlashDateWeekdayOnly,
		},
		// パターン8: "10/6 18:00まで" (スラッシュ日付+時刻まで)
		{
			regexp.MustCompile(`(\d{1,2})/(\d{1,2})[\s　]+(\d{1,2}):(\d{2})\s*まで`),
			s.handleSlashDateTimeUntil,
		},
		// パターン9: "10/6 18:00-20:00" (スラッシュ日付+時刻範囲)
		{
			regexp.MustCompile(`(\d{1,2})/(\d{1,2})[\s　]+(\d{1,2}):(\d{2})\s*[-〜～]\s*(\d{1,2}):(\d{2})`),
			s.handleSlashDateTimeRange,
		},
		// パターン10: "10/6 18:00" (スラッシュ日付+時刻)
		{
			regexp.MustCompile(`(\d{1,2})/(\d{1,2})[\s　]+(\d{1,2}):(\d{2})`),
			s.handleSlashDateTime,
		},
		// パターン11: "10/6" (スラッシュ日付のみ)
		{
			regexp.MustCompile(`(\d{1,2})/(\d{1,2})`),
			s.handleSlashDateOnly,
		},
		// パターン12: "10月3日（月）18:00まで" (月日曜日+時刻まで)
		{
			regexp.MustCompile(`(\d{1,2})月(\d{1,2})日（[月火水木金土日]）[\s　]*(\d{1,2}):(\d{2})\s*まで`),
			s.handleDateWeekdayTimeUntil,
		},
		// パターン13: "10月3日（月）18:00-20:00" (月日曜日+時刻範囲)
		{
			regexp.MustCompile(`(\d{1,2})月(\d{1,2})日（[月火水木金土日]）[\s　]*(\d{1,2}):(\d{2})\s*[-〜～]\s*(\d{1,2}):(\d{2})`),
			s.handleDateWeekdayTimeRange,
		},
		// パターン14: "10月3日（月）18:00" (月日曜日+時刻)
		{
			regexp.MustCompile(`(\d{1,2})月(\d{1,2})日（[月火水木金土日]）[\s　]*(\d{1,2}):(\d{2})`),
			s.handleDateWeekdayTime,
		},
		// パターン15: "10月3日（月）" (月日曜日のみ)
		{
			regexp.MustCompile(`(\d{1,2})月(\d{1,2})日（[月火水木金土日]）`),
			s.handleDateWeekdayOnly,
		},
		// パターン16: "10月3日 14:00-16:00" (月日+時刻範囲)
		{
			regexp.MustCompile(`(\d{1,2})月(\d{1,2})日[\s　]*(\d{1,2}):(\d{2})\s*[-〜～]\s*(\d{1,2}):(\d{2})`),
			s.handleDateTimeRange,
		},
		// パターン17: "10月3日 14:00" (月日+時刻)
		{
			regexp.MustCompile(`(\d{1,2})月(\d{1,2})日[\s　]*(\d{1,2}):(\d{2})`),
			s.handleDateTime,
		},
		// パターン18: "14:00-16:00" (時刻範囲のみ)
		{
			regexp.MustCompile(`(\d{1,2}):(\d{2})\s*[-〜～]\s*(\d{1,2}):(\d{2})`),
			s.handleTimeRange,
		},
		// パターン19: "14時30分〜16時45分" (時分範囲・日本語)
		{
			regexp.MustCompile(`(\d{1,2})時(\d{1,2})分\s*[〜～]\s*(\d{1,2})時(\d{1,2})分`),
			s.handleJapaneseTimeMinuteRange,
		},
		// パターン20: "14時30分から16時45分" (時分範囲・から)
		{
			regexp.MustCompile(`(\d{1,2})時(\d{1,2})分から\s*(\d{1,2})時(\d{1,2})分`),
			s.handleJapaneseTimeMinuteFromTo,
		},
		// パターン21: "14時30分から" (時分開始・から)
		{
			regexp.MustCompile(`(\d{1,2})時(\d{1,2})分から[！!]?`),
			s.handleJapaneseTimeMinuteFrom,
		},
		// パターン22: "14時30分〜" (時分開始・〜)
		{
			regexp.MustCompile(`(\d{1,2})時(\d{1,2})分[〜～][！!]?`),
			s.handleJapaneseTimeMinuteStart,
		},
		// パターン23: "14時30分" (時分のみ)
		{
			regexp.MustCompile(`(\d{1,2})時(\d{1,2})分`),
			s.handleJapaneseTimeMinute,
		},
		// パターン24: "14時〜16時" (時刻範囲・日本語)
		{
			regexp.MustCompile(`(\d{1,2})時\s*[〜～]\s*(\d{1,2})時`),
			s.handleJapaneseTimeRange,
		},
		// パターン25: "14時から16時" (時刻範囲・から)
		{
			regexp.MustCompile(`(\d{1,2})時から\s*(\d{1,2})時`),
			s.handleJapaneseTimeFromTo,
		},
		// パターン26: "14時から" (開始時刻のみ・から)
		{
			regexp.MustCompile(`(\d{1,2})時から[！!]?`),
			s.handleJapaneseTimeFrom,
		},
		// パターン27: "14時〜" (開始時刻のみ・〜)
		{
			regexp.MustCompile(`(\d{1,2})時[〜～][！!]?`),
			s.handleJapaneseTimeStart,
		},
		// === 相対日付表現 ===（具体的なパターンを先に配置）
		// パターン28: "明日 14:00"
		{
			regexp.MustCompile(`明日[\s　]*(\d{1,2}):(\d{2})`),
			s.handleTomorrowTime,
		},
		// パターン29: "今日 14:00"
		{
			regexp.MustCompile(`今日[\s　]*(\d{1,2}):(\d{2})`),
			s.handleTodayTime,
		},
		// パターン30: "明後日 14:00"
		{
			regexp.MustCompile(`明後日[\s　]*(\d{1,2}):(\d{2})`),
			s.handleDayAfterTomorrowTime,
		},

		// === 英語混在表現 ===（時刻のみより先に配置）
		// パターン31: "AM 9:00", "PM 6:00"
		{
			regexp.MustCompile(`AM[\s　]*(\d{1,2}):(\d{2})`),
			s.handleAMTimeEng,
		},
		{
			regexp.MustCompile(`PM[\s　]*(\d{1,2}):(\d{2})`),
			s.handlePMTimeEng,
		},

		// === 区切り文字バリエーション ===（時刻のみより先に配置）
		// パターン32: "10-6 18:00", "10.6 18:00"
		{
			regexp.MustCompile(`(\d{1,2})[-.](\d{1,2})[\s　]+(\d{1,2}):(\d{2})`),
			s.handleAlternativeDateFormat,
		},

		// === 自然な日本語表現 ===（時刻のみパターンより先に配置）
		// パターン33-42: "今週の土曜日", "今度の土曜日", "来週月曜日"
		{
			regexp.MustCompile(`今週の?([月火水木金土日])曜日[\s　]*(\d{1,2}):(\d{2})`),
			s.handleThisWeekdayTime,
		},
		{
			regexp.MustCompile(`今度の([月火水木金土日])曜日[\s　]*(\d{1,2}):(\d{2})`),
			s.handleNextWeekdayTime,
		},
		{
			regexp.MustCompile(`来週の?([月火水木金土日])曜日[\s　]*(\d{1,2}):(\d{2})`),
			s.handleNextWeekWeekdayTime,
		},
		{
			regexp.MustCompile(`今週の?([月火水木金土日])曜日`),
			s.handleThisWeekday,
		},
		{
			regexp.MustCompile(`今度の([月火水木金土日])曜日`),
			s.handleNextWeekday,
		},
		{
			regexp.MustCompile(`来週の?([月火水木金土日])曜日`),
			s.handleNextWeekWeekday,
		},

		// パターン37: "14:00" (時刻のみ) - より具体的なパターンの後に配置
		{
			regexp.MustCompile(`(\d{1,2}):(\d{2})`),
			s.handleTimeOnly,
		},
		// パターン38: "10月3日" (月日のみ)
		{
			regexp.MustCompile(`(\d{1,2})月(\d{1,2})日`),
			s.handleDateOnly,
		},

		// === 時間帯表現 ===
		// パターン39: "午前10時", "午後3時"
		{
			regexp.MustCompile(`午前(\d{1,2})時`),
			s.handleAMTime,
		},
		{
			regexp.MustCompile(`午後(\d{1,2})時`),
			s.handlePMTime,
		},
		// パターン40: "夜8時", "朝9時", "昼12時"
		{
			regexp.MustCompile(`夜(\d{1,2})時`),
			s.handleNightTime,
		},
		{
			regexp.MustCompile(`朝(\d{1,2})時`),
			s.handleMorningTime,
		},
		{
			regexp.MustCompile(`昼(\d{1,2})時`),
			s.handleNoonTime,
		},

		// === 完全日付形式 ===
		// パターン41: "2025/10/6 18:00"
		{
			regexp.MustCompile(`(\d{4})/(\d{1,2})/(\d{1,2})[\s　]*(\d{1,2}):(\d{2})`),
			s.handleFullSlashDate,
		},

		// === 曖昧な時間表現 ===
		// パターン42-45: "夕方", "お昼頃", "夜中", "早朝"
		{
			regexp.MustCompile(`夕方`),
			s.handleEvening,
		},
		{
			regexp.MustCompile(`お昼頃|昼頃`),
			s.handleAroundNoon,
		},
		{
			regexp.MustCompile(`夜中|深夜`),
			s.handleMidnight,
		},
		{
			regexp.MustCompile(`早朝`),
			s.handleEarlyMorning,
		},

		// === 期間表現 ===
		// パターン46-47: "3時間", "30分間"
		{
			regexp.MustCompile(`(\d{1,2})時間`),
			s.handleHourDuration,
		},
		{
			regexp.MustCompile(`(\d{1,2})分間`),
			s.handleMinuteDuration,
		},
	}

	// 各パターンを試行
	for _, pattern := range patterns {
		matches := pattern.regex.FindStringSubmatch(content)
		if len(matches) > 0 {
			log.Printf("DateTime extraction - Pattern matched: %v", matches)
			startsAt, endsAt := pattern.handler(matches, postCreatedAt)
			return startsAt, endsAt, true
		}
	}

	// パターンが見つからない場合はデフォルト（投稿日の0:00-1:00）を返すが、パターンマッチしなかったことを示す
	log.Printf("DateTime extraction - No pattern found, using default time")
	startsAt, endsAt := s.getDefaultDateTime(postCreatedAt)
	return startsAt, endsAt, false
}

// 年月日+時刻範囲の処理 (例: "2026年1月10日 14:00-16:00")
func (s *DateTimeExtractionService) handleYearDateTimeRange(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
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
func (s *DateTimeExtractionService) handleYearDateTime(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
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
func (s *DateTimeExtractionService) handleYearDateOnly(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	year, _ := strconv.Atoi(matches[1])
	month, _ := strconv.Atoi(matches[2])
	day, _ := strconv.Atoi(matches[3])

	startsAt := time.Date(year, time.Month(month), day, 0, 0, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// 時刻範囲の処理 (例: "14:00-16:00")
func (s *DateTimeExtractionService) handleTimeRange(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
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

// 時刻のみの処理 (例: "14:00")
func (s *DateTimeExtractionService) handleTimeOnly(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	hour, _ := strconv.Atoi(matches[1])
	min, _ := strconv.Atoi(matches[2])

	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		hour, min, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// 日付のみの処理 (例: "10月3日")
func (s *DateTimeExtractionService) handleDateOnly(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
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

// 日付+時刻の処理 (例: "10月3日 14:00")
func (s *DateTimeExtractionService) handleDateTime(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	month, _ := strconv.Atoi(matches[1])
	day, _ := strconv.Atoi(matches[2])
	hour, _ := strconv.Atoi(matches[3])
	min, _ := strconv.Atoi(matches[4])

	year := postCreatedAt.Year()
	if month < int(postCreatedAt.Month()) {
		year++
	}

	startsAt := time.Date(year, time.Month(month), day, hour, min, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// 日付+時刻範囲の処理 (例: "10月3日 14:00-16:00")
func (s *DateTimeExtractionService) handleDateTimeRange(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	month, _ := strconv.Atoi(matches[1])
	day, _ := strconv.Atoi(matches[2])
	startHour, _ := strconv.Atoi(matches[3])
	startMin, _ := strconv.Atoi(matches[4])
	endHour, _ := strconv.Atoi(matches[5])
	endMin, _ := strconv.Atoi(matches[6])

	year := postCreatedAt.Year()
	if month < int(postCreatedAt.Month()) {
		year++
	}

	startsAt := time.Date(year, time.Month(month), day, startHour, startMin, 0, 0, postCreatedAt.Location())
	endsAt := time.Date(year, time.Month(month), day, endHour, endMin, 0, 0, postCreatedAt.Location())

	return startsAt, &endsAt
}

// 簡略化のため、他のハンドラー関数は基本的なものを実装
// 実際の本格運用時には全パターンを実装する必要があります

// スラッシュ形式の基本ハンドラー（簡略版）
func (s *DateTimeExtractionService) handleSlashDateWeekdayTimeUntil(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	month, _ := strconv.Atoi(matches[1])
	day, _ := strconv.Atoi(matches[2])
	hour, _ := strconv.Atoi(matches[3])
	min, _ := strconv.Atoi(matches[4])

	year := postCreatedAt.Year()
	if month < int(postCreatedAt.Month()) {
		year++
	}

	// "まで"の場合は終了時刻として扱い、開始時刻は投稿時刻とする
	startsAt := postCreatedAt
	endsAt := time.Date(year, time.Month(month), day, hour, min, 0, 0, postCreatedAt.Location())

	return startsAt, &endsAt
}

func (s *DateTimeExtractionService) handleSlashDateWeekdayTimeRange(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	return s.handleDateTimeRange(matches, postCreatedAt)
}

func (s *DateTimeExtractionService) handleSlashDateWeekdayTime(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	return s.handleDateTime(matches, postCreatedAt)
}

func (s *DateTimeExtractionService) handleSlashDateWeekdayOnly(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	return s.handleDateOnly(matches, postCreatedAt)
}

func (s *DateTimeExtractionService) handleDateWeekdayTimeUntil(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	month, _ := strconv.Atoi(matches[1])
	day, _ := strconv.Atoi(matches[2])
	hour, _ := strconv.Atoi(matches[3])
	min, _ := strconv.Atoi(matches[4])

	year := postCreatedAt.Year()
	if month < int(postCreatedAt.Month()) {
		year++
	}

	// "まで"の場合は終了時刻として扱い、開始時刻は投稿時刻とする
	startsAt := postCreatedAt
	endsAt := time.Date(year, time.Month(month), day, hour, min, 0, 0, postCreatedAt.Location())

	return startsAt, &endsAt
}

func (s *DateTimeExtractionService) handleDateWeekdayTimeRange(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	return s.handleDateTimeRange(matches, postCreatedAt)
}

func (s *DateTimeExtractionService) handleDateWeekdayTime(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	return s.handleDateTime(matches, postCreatedAt)
}

func (s *DateTimeExtractionService) handleDateWeekdayOnly(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	return s.handleDateOnly(matches, postCreatedAt)
}

// 日本語時刻パターンの基本ハンドラー（簡略版）
func (s *DateTimeExtractionService) handleJapaneseTimeMinuteRange(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	return s.handleTimeRange(matches, postCreatedAt)
}

func (s *DateTimeExtractionService) handleJapaneseTimeMinuteFromTo(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	return s.handleTimeRange(matches, postCreatedAt)
}

func (s *DateTimeExtractionService) handleJapaneseTimeMinuteFrom(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	return s.handleTimeOnly(matches, postCreatedAt)
}

func (s *DateTimeExtractionService) handleJapaneseTimeMinuteStart(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	return s.handleTimeOnly(matches, postCreatedAt)
}

func (s *DateTimeExtractionService) handleJapaneseTimeMinute(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	return s.handleTimeOnly(matches, postCreatedAt)
}

func (s *DateTimeExtractionService) handleJapaneseTimeRange(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	startHour, _ := strconv.Atoi(matches[1])
	endHour, _ := strconv.Atoi(matches[2])

	// 投稿日の指定時刻（分は0とする）
	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		startHour, 0, 0, 0, postCreatedAt.Location())
	endsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		endHour, 0, 0, 0, postCreatedAt.Location())

	return startsAt, &endsAt
}

func (s *DateTimeExtractionService) handleJapaneseTimeFromTo(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	startHour, _ := strconv.Atoi(matches[1])
	endHour, _ := strconv.Atoi(matches[2])

	// 投稿日の指定時刻（分は0とする）
	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		startHour, 0, 0, 0, postCreatedAt.Location())
	endsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		endHour, 0, 0, 0, postCreatedAt.Location())

	return startsAt, &endsAt
}

func (s *DateTimeExtractionService) handleJapaneseTimeFrom(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	hour, _ := strconv.Atoi(matches[1])

	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		hour, 0, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour) // 1時間後を終了時刻とする

	return startsAt, &endsAt
}

func (s *DateTimeExtractionService) handleJapaneseTimeStart(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	hour, _ := strconv.Atoi(matches[1])

	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		hour, 0, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour) // 1時間後を終了時刻とする

	return startsAt, &endsAt
}

// スラッシュ日付+時刻まで の処理 (例: "10/6 18:00まで")
func (s *DateTimeExtractionService) handleSlashDateTimeUntil(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	month, _ := strconv.Atoi(matches[1])
	day, _ := strconv.Atoi(matches[2])
	hour, _ := strconv.Atoi(matches[3])
	min, _ := strconv.Atoi(matches[4])

	year := postCreatedAt.Year()
	if month < int(postCreatedAt.Month()) {
		year++
	}

	// "まで"の場合は終了時刻として扱い、開始時刻は投稿時刻とする
	startsAt := postCreatedAt
	endsAt := time.Date(year, time.Month(month), day, hour, min, 0, 0, postCreatedAt.Location())

	return startsAt, &endsAt
}

// スラッシュ日付+時刻範囲 の処理 (例: "10/6 18:00-20:00")
func (s *DateTimeExtractionService) handleSlashDateTimeRange(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	month, _ := strconv.Atoi(matches[1])
	day, _ := strconv.Atoi(matches[2])
	startHour, _ := strconv.Atoi(matches[3])
	startMin, _ := strconv.Atoi(matches[4])
	endHour, _ := strconv.Atoi(matches[5])
	endMin, _ := strconv.Atoi(matches[6])

	year := postCreatedAt.Year()
	if month < int(postCreatedAt.Month()) {
		year++
	}

	startsAt := time.Date(year, time.Month(month), day, startHour, startMin, 0, 0, postCreatedAt.Location())
	endsAt := time.Date(year, time.Month(month), day, endHour, endMin, 0, 0, postCreatedAt.Location())

	return startsAt, &endsAt
}

// スラッシュ日付+時刻 の処理 (例: "10/6 18:00")
func (s *DateTimeExtractionService) handleSlashDateTime(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	month, _ := strconv.Atoi(matches[1])
	day, _ := strconv.Atoi(matches[2])
	hour, _ := strconv.Atoi(matches[3])
	min, _ := strconv.Atoi(matches[4])

	year := postCreatedAt.Year()
	if month < int(postCreatedAt.Month()) {
		year++
	}

	startsAt := time.Date(year, time.Month(month), day, hour, min, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// スラッシュ日付のみ の処理 (例: "10/6")
func (s *DateTimeExtractionService) handleSlashDateOnly(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
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

// === 相対日付表現ハンドラー ===

// 明日の時刻処理 (例: "明日 14:00")
func (s *DateTimeExtractionService) handleTomorrowTime(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	hour, _ := strconv.Atoi(matches[1])
	min, _ := strconv.Atoi(matches[2])

	tomorrow := postCreatedAt.AddDate(0, 0, 1)
	startsAt := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(),
		hour, min, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// 今日の時刻処理 (例: "今日 14:00")
func (s *DateTimeExtractionService) handleTodayTime(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	hour, _ := strconv.Atoi(matches[1])
	min, _ := strconv.Atoi(matches[2])

	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		hour, min, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// 明後日の時刻処理 (例: "明後日 14:00")
func (s *DateTimeExtractionService) handleDayAfterTomorrowTime(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	hour, _ := strconv.Atoi(matches[1])
	min, _ := strconv.Atoi(matches[2])

	dayAfterTomorrow := postCreatedAt.AddDate(0, 0, 2)
	startsAt := time.Date(dayAfterTomorrow.Year(), dayAfterTomorrow.Month(), dayAfterTomorrow.Day(),
		hour, min, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// === 時間帯表現ハンドラー ===

// 午前時刻処理 (例: "午前10時")
func (s *DateTimeExtractionService) handleAMTime(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	hour, _ := strconv.Atoi(matches[1])
	if hour == 12 {
		hour = 0 // 午前12時は0時
	}

	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		hour, 0, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// 午後時刻処理 (例: "午後3時")
func (s *DateTimeExtractionService) handlePMTime(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	hour, _ := strconv.Atoi(matches[1])
	if hour != 12 {
		hour += 12 // 午後は12時間追加（12時は除く）
	}

	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		hour, 0, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// 夜時刻処理 (例: "夜8時")
func (s *DateTimeExtractionService) handleNightTime(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	hour, _ := strconv.Atoi(matches[1])
	if hour < 12 {
		hour += 12 // 夜は午後として扱う
	}

	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		hour, 0, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// 朝時刻処理 (例: "朝9時")
func (s *DateTimeExtractionService) handleMorningTime(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	hour, _ := strconv.Atoi(matches[1])

	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		hour, 0, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// 昼時刻処理 (例: "昼12時")
func (s *DateTimeExtractionService) handleNoonTime(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	hour, _ := strconv.Atoi(matches[1])

	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		hour, 0, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// === 英語混在表現ハンドラー ===

// AM時刻処理 (例: "AM 9:00")
func (s *DateTimeExtractionService) handleAMTimeEng(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	hour, _ := strconv.Atoi(matches[1])
	min, _ := strconv.Atoi(matches[2])
	if hour == 12 {
		hour = 0 // AM 12:00は0:00
	}

	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		hour, min, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// PM時刻処理 (例: "PM 6:00")
func (s *DateTimeExtractionService) handlePMTimeEng(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	hour, _ := strconv.Atoi(matches[1])
	min, _ := strconv.Atoi(matches[2])
	if hour != 12 {
		hour += 12 // PM時は12時間追加（12時は除く）
	}

	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		hour, min, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// === 区切り文字バリエーションハンドラー ===

// 代替日付形式処理 (例: "10-6 18:00", "10.6 18:00")
func (s *DateTimeExtractionService) handleAlternativeDateFormat(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	month, _ := strconv.Atoi(matches[1])
	day, _ := strconv.Atoi(matches[2])
	hour, _ := strconv.Atoi(matches[3])
	min, _ := strconv.Atoi(matches[4])

	year := postCreatedAt.Year()
	if month < int(postCreatedAt.Month()) {
		year++
	}

	startsAt := time.Date(year, time.Month(month), day, hour, min, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// 完全スラッシュ日付処理 (例: "2025/10/6 18:00")
func (s *DateTimeExtractionService) handleFullSlashDate(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	year, _ := strconv.Atoi(matches[1])
	month, _ := strconv.Atoi(matches[2])
	day, _ := strconv.Atoi(matches[3])
	hour, _ := strconv.Atoi(matches[4])
	min, _ := strconv.Atoi(matches[5])

	startsAt := time.Date(year, time.Month(month), day, hour, min, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// === 曖昧な時間表現ハンドラー ===

// 夕方処理 (例: "夕方")
func (s *DateTimeExtractionService) handleEvening(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		18, 0, 0, 0, postCreatedAt.Location()) // 18:00を夕方とする
	endsAt := startsAt.Add(2 * time.Hour) // 2時間の幅

	return startsAt, &endsAt
}

// お昼頃処理 (例: "お昼頃")
func (s *DateTimeExtractionService) handleAroundNoon(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		12, 0, 0, 0, postCreatedAt.Location()) // 12:00をお昼とする
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// 夜中・深夜処理 (例: "夜中", "深夜")
func (s *DateTimeExtractionService) handleMidnight(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		0, 0, 0, 0, postCreatedAt.Location()) // 0:00を夜中とする
	endsAt := startsAt.Add(2 * time.Hour) // 2時間の幅

	return startsAt, &endsAt
}

// 早朝処理 (例: "早朝")
func (s *DateTimeExtractionService) handleEarlyMorning(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		6, 0, 0, 0, postCreatedAt.Location()) // 6:00を早朝とする
	endsAt := startsAt.Add(2 * time.Hour) // 2時間の幅

	return startsAt, &endsAt
}

// === 自然な日本語表現ハンドラー ===

// 今週の曜日+時刻処理 (例: "今週の土曜日 14:00")
func (s *DateTimeExtractionService) handleThisWeekdayTime(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	weekdayStr := matches[1]
	hour, _ := strconv.Atoi(matches[2])
	min, _ := strconv.Atoi(matches[3])

	targetWeekday := s.parseWeekday(weekdayStr)
	daysUntil := s.daysUntilWeekday(postCreatedAt.Weekday(), targetWeekday)

	// 今週の場合、過去の曜日も含める（例：金曜日に「今週の月曜日」と言った場合、今週の月曜日を指す）
	if daysUntil > 0 && targetWeekday < postCreatedAt.Weekday() {
		// 対象曜日が今日より前の曜日の場合、今週のその曜日（過去）を指す
		daysUntil -= 7
	}

	targetDate := postCreatedAt.AddDate(0, 0, daysUntil)
	startsAt := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(),
		hour, min, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// 今週の曜日処理 (例: "今週の土曜日")
func (s *DateTimeExtractionService) handleThisWeekday(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	weekdayStr := matches[1]

	targetWeekday := s.parseWeekday(weekdayStr)
	daysUntil := s.daysUntilWeekday(postCreatedAt.Weekday(), targetWeekday)

	// 今週の場合、過去の曜日も含める
	if daysUntil > 0 && targetWeekday < postCreatedAt.Weekday() {
		// 対象曜日が今日より前の曜日の場合、今週のその曜日（過去）を指す
		daysUntil -= 7
	}

	targetDate := postCreatedAt.AddDate(0, 0, daysUntil)
	startsAt := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(),
		0, 0, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// 今度の曜日+時刻処理 (例: "今度の土曜日 14:00")
func (s *DateTimeExtractionService) handleNextWeekdayTime(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	weekdayStr := matches[1]
	hour, _ := strconv.Atoi(matches[2])
	min, _ := strconv.Atoi(matches[3])

	targetWeekday := s.parseWeekday(weekdayStr)
	daysUntil := s.daysUntilWeekday(postCreatedAt.Weekday(), targetWeekday)
	if daysUntil == 0 {
		daysUntil = 7 // 今日が同じ曜日なら来週
	}

	targetDate := postCreatedAt.AddDate(0, 0, daysUntil)
	startsAt := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(),
		hour, min, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// 来週曜日+時刻処理 (例: "来週月曜日 10:00")
func (s *DateTimeExtractionService) handleNextWeekWeekdayTime(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	weekdayStr := matches[1]
	hour, _ := strconv.Atoi(matches[2])
	min, _ := strconv.Atoi(matches[3])

	targetWeekday := s.parseWeekday(weekdayStr)
	daysUntil := s.daysUntilWeekday(postCreatedAt.Weekday(), targetWeekday) + 7 // 来週なので+7日

	targetDate := postCreatedAt.AddDate(0, 0, daysUntil)
	startsAt := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(),
		hour, min, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// 今度の曜日処理 (例: "今度の土曜日")
func (s *DateTimeExtractionService) handleNextWeekday(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	weekdayStr := matches[1]

	targetWeekday := s.parseWeekday(weekdayStr)
	daysUntil := s.daysUntilWeekday(postCreatedAt.Weekday(), targetWeekday)
	if daysUntil == 0 {
		daysUntil = 7 // 今日が同じ曜日なら来週
	}

	targetDate := postCreatedAt.AddDate(0, 0, daysUntil)
	startsAt := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(),
		0, 0, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// 来週曜日処理 (例: "来週月曜日")
func (s *DateTimeExtractionService) handleNextWeekWeekday(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	weekdayStr := matches[1]

	targetWeekday := s.parseWeekday(weekdayStr)
	daysUntil := s.daysUntilWeekday(postCreatedAt.Weekday(), targetWeekday) + 7 // 来週なので+7日

	targetDate := postCreatedAt.AddDate(0, 0, daysUntil)
	startsAt := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(),
		0, 0, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

// === 期間表現ハンドラー ===

// 時間期間処理 (例: "3時間")
func (s *DateTimeExtractionService) handleHourDuration(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	hours, _ := strconv.Atoi(matches[1])

	startsAt := postCreatedAt
	endsAt := startsAt.Add(time.Duration(hours) * time.Hour)

	return startsAt, &endsAt
}

// 分期間処理 (例: "30分間")
func (s *DateTimeExtractionService) handleMinuteDuration(matches []string, postCreatedAt time.Time) (time.Time, *time.Time) {
	minutes, _ := strconv.Atoi(matches[1])

	startsAt := postCreatedAt
	endsAt := startsAt.Add(time.Duration(minutes) * time.Minute)

	return startsAt, &endsAt
}

// === ヘルパー関数 ===

// 曜日文字列を time.Weekday に変換
func (s *DateTimeExtractionService) parseWeekday(weekdayStr string) time.Weekday {
	weekdayMap := map[string]time.Weekday{
		"日": time.Sunday,
		"月": time.Monday,
		"火": time.Tuesday,
		"水": time.Wednesday,
		"木": time.Thursday,
		"金": time.Friday,
		"土": time.Saturday,
	}
	return weekdayMap[weekdayStr]
}

// 指定曜日まで何日かを計算
func (s *DateTimeExtractionService) daysUntilWeekday(current, target time.Weekday) int {
	days := int(target - current)
	if days < 0 {
		days += 7
	}
	return days
}

// デフォルト日時の取得 (投稿日の0:00-1:00)
func (s *DateTimeExtractionService) getDefaultDateTime(postCreatedAt time.Time) (time.Time, *time.Time) {
	startsAt := time.Date(postCreatedAt.Year(), postCreatedAt.Month(), postCreatedAt.Day(),
		0, 0, 0, 0, postCreatedAt.Location())
	endsAt := startsAt.Add(1 * time.Hour)

	return startsAt, &endsAt
}

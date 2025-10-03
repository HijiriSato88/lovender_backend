package service

import (
	"testing"
	"time"
)

// time.Timeのポインタを返すヘルパー関数
func timePtr(t time.Time) *time.Time {
	return &t
}

func TestDateTimeExtractionService_ExtractDateTime_Complete109Patterns(t *testing.T) {
	service := NewDateTimeExtractionService()

	// テスト用の基準日時（2025年10月3日 12:00:00）
	baseTime := time.Date(2025, 10, 3, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		no            int
		name          string
		content       string
		expectedStart time.Time
		expectedEnd   *time.Time
		description   string
	}{
		// 1-7: 年月日指定パターン
		{
			no:            1,
			name:          "年月日+時刻範囲（スペースなし）",
			content:       "2026年1月10日14:00-16:00のイベント",
			expectedStart: time.Date(2026, 1, 10, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2026, 1, 10, 16, 0, 0, 0, time.UTC)),
			description:   "年月日+時刻範囲、スペースなし",
		},
		{
			no:            2,
			name:          "年月日+時刻範囲（半角スペース）",
			content:       "2026年1月10日 14:00-16:00のイベント",
			expectedStart: time.Date(2026, 1, 10, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2026, 1, 10, 16, 0, 0, 0, time.UTC)),
			description:   "年月日+時刻範囲、半角スペース",
		},
		{
			no:            3,
			name:          "年月日+時刻範囲（全角スペース）",
			content:       "2026年1月10日　14:00-16:00のイベント",
			expectedStart: time.Date(2026, 1, 10, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2026, 1, 10, 16, 0, 0, 0, time.UTC)),
			description:   "年月日+時刻範囲、全角スペース",
		},
		{
			no:            4,
			name:          "年月日+時刻（スペースなし）",
			content:       "2026年1月10日14:00開始",
			expectedStart: time.Date(2026, 1, 10, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2026, 1, 10, 15, 0, 0, 0, time.UTC)),
			description:   "年月日+時刻、スペースなし",
		},
		{
			no:            5,
			name:          "年月日+時刻（半角スペース）",
			content:       "2026年1月10日 14:00開始",
			expectedStart: time.Date(2026, 1, 10, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2026, 1, 10, 15, 0, 0, 0, time.UTC)),
			description:   "年月日+時刻、半角スペース",
		},
		{
			no:            6,
			name:          "年月日+時刻（全角スペース）",
			content:       "2026年1月10日　14:00開始",
			expectedStart: time.Date(2026, 1, 10, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2026, 1, 10, 15, 0, 0, 0, time.UTC)),
			description:   "年月日+時刻、全角スペース",
		},
		{
			no:            7,
			name:          "年月日のみ",
			content:       "2026年1月10日のイベント",
			expectedStart: time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2026, 1, 10, 1, 0, 0, 0, time.UTC)),
			description:   "年月日のみ",
		},

		// 8-17: スラッシュ形式（曜日付き）
		{
			no:            8,
			name:          "曜日+時刻まで（スペースなし）",
			content:       "10/6（月）18:00まで",
			expectedStart: baseTime,
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 18, 0, 0, 0, time.UTC)),
			description:   "スラッシュ曜日+時刻まで、スペースなし",
		},
		{
			no:            9,
			name:          "曜日+時刻まで（半角スペース）",
			content:       "10/6（月） 18:00まで",
			expectedStart: baseTime,
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 18, 0, 0, 0, time.UTC)),
			description:   "スラッシュ曜日+時刻まで、半角スペース",
		},
		{
			no:            10,
			name:          "曜日+時刻まで（全角スペース）",
			content:       "10/6（月）　18:00まで",
			expectedStart: baseTime,
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 18, 0, 0, 0, time.UTC)),
			description:   "スラッシュ曜日+時刻まで、全角スペース",
		},
		{
			no:            11,
			name:          "曜日+時刻範囲（スペースなし）",
			content:       "10/6（月）18:00-20:00",
			expectedStart: time.Date(2025, 10, 6, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 20, 0, 0, 0, time.UTC)),
			description:   "スラッシュ曜日+時刻範囲、スペースなし",
		},
		{
			no:            12,
			name:          "曜日+時刻範囲（半角スペース）",
			content:       "10/6（月） 18:00-20:00",
			expectedStart: time.Date(2025, 10, 6, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 20, 0, 0, 0, time.UTC)),
			description:   "スラッシュ曜日+時刻範囲、半角スペース",
		},
		{
			no:            13,
			name:          "曜日+時刻範囲（全角スペース）",
			content:       "10/6（月）　18:00-20:00",
			expectedStart: time.Date(2025, 10, 6, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 20, 0, 0, 0, time.UTC)),
			description:   "スラッシュ曜日+時刻範囲、全角スペース",
		},
		{
			no:            14,
			name:          "曜日+時刻（スペースなし）",
			content:       "10/6（月）18:00",
			expectedStart: time.Date(2025, 10, 6, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 19, 0, 0, 0, time.UTC)),
			description:   "スラッシュ曜日+時刻、スペースなし",
		},
		{
			no:            15,
			name:          "曜日+時刻（半角スペース）",
			content:       "10/6（月） 18:00",
			expectedStart: time.Date(2025, 10, 6, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 19, 0, 0, 0, time.UTC)),
			description:   "スラッシュ曜日+時刻、半角スペース",
		},
		{
			no:            16,
			name:          "曜日+時刻（全角スペース）",
			content:       "10/6（月）　18:00",
			expectedStart: time.Date(2025, 10, 6, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 19, 0, 0, 0, time.UTC)),
			description:   "スラッシュ曜日+時刻、全角スペース",
		},
		{
			no:            17,
			name:          "曜日のみ",
			content:       "10/6（月）",
			expectedStart: time.Date(2025, 10, 6, 0, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 1, 0, 0, 0, time.UTC)),
			description:   "スラッシュ曜日のみ",
		},

		// 18-24: スラッシュ形式（曜日なし）
		{
			no:            18,
			name:          "日付+時刻まで（半角スペース）",
			content:       "10/6 18:00まで",
			expectedStart: baseTime,
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 18, 0, 0, 0, time.UTC)),
			description:   "スラッシュ日付+時刻まで、半角スペース",
		},
		{
			no:            19,
			name:          "日付+時刻まで（全角スペース）",
			content:       "10/6　18:00まで",
			expectedStart: baseTime,
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 18, 0, 0, 0, time.UTC)),
			description:   "スラッシュ日付+時刻まで、全角スペース",
		},
		{
			no:            20,
			name:          "日付+時刻範囲（半角スペース）",
			content:       "10/6 18:00-20:00",
			expectedStart: time.Date(2025, 10, 6, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 20, 0, 0, 0, time.UTC)),
			description:   "スラッシュ日付+時刻範囲、半角スペース",
		},
		{
			no:            21,
			name:          "日付+時刻範囲（全角スペース）",
			content:       "10/6　18:00-20:00",
			expectedStart: time.Date(2025, 10, 6, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 20, 0, 0, 0, time.UTC)),
			description:   "スラッシュ日付+時刻範囲、全角スペース",
		},
		{
			no:            22,
			name:          "日付+時刻（半角スペース）",
			content:       "10/6 18:00",
			expectedStart: time.Date(2025, 10, 6, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 19, 0, 0, 0, time.UTC)),
			description:   "スラッシュ日付+時刻、半角スペース",
		},
		{
			no:            23,
			name:          "日付+時刻（全角スペース）",
			content:       "10/6　18:00",
			expectedStart: time.Date(2025, 10, 6, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 19, 0, 0, 0, time.UTC)),
			description:   "スラッシュ日付+時刻、全角スペース",
		},
		{
			no:            24,
			name:          "日付のみ",
			content:       "10/6",
			expectedStart: time.Date(2025, 10, 6, 0, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 1, 0, 0, 0, time.UTC)),
			description:   "スラッシュ日付のみ",
		},

		// 25-33: 月日形式（曜日付き）
		{
			no:            25,
			name:          "月日曜日+時刻まで（スペースなし）",
			content:       "10月3日（月）18:00まで",
			expectedStart: baseTime,
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 18, 0, 0, 0, time.UTC)),
			description:   "月日曜日+時刻まで、スペースなし",
		},
		{
			no:            26,
			name:          "月日曜日+時刻まで（半角スペース）",
			content:       "10月3日（月） 18:00まで",
			expectedStart: baseTime,
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 18, 0, 0, 0, time.UTC)),
			description:   "月日曜日+時刻まで、半角スペース",
		},
		{
			no:            27,
			name:          "月日曜日+時刻まで（全角スペース）",
			content:       "10月3日（月）　18:00まで",
			expectedStart: baseTime,
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 18, 0, 0, 0, time.UTC)),
			description:   "月日曜日+時刻まで、全角スペース",
		},
		{
			no:            28,
			name:          "月日曜日+時刻範囲（スペースなし）",
			content:       "10月3日（月）18:00-20:00",
			expectedStart: time.Date(2025, 10, 3, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 20, 0, 0, 0, time.UTC)),
			description:   "月日曜日+時刻範囲、スペースなし",
		},
		{
			no:            29,
			name:          "月日曜日+時刻範囲（半角スペース）",
			content:       "10月3日（月） 18:00-20:00",
			expectedStart: time.Date(2025, 10, 3, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 20, 0, 0, 0, time.UTC)),
			description:   "月日曜日+時刻範囲、半角スペース",
		},
		{
			no:            30,
			name:          "月日曜日+時刻範囲（全角スペース）",
			content:       "10月3日（月）　18:00-20:00",
			expectedStart: time.Date(2025, 10, 3, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 20, 0, 0, 0, time.UTC)),
			description:   "月日曜日+時刻範囲、全角スペース",
		},
		{
			no:            31,
			name:          "月日曜日+時刻（スペースなし）",
			content:       "10月3日（月）18:00",
			expectedStart: time.Date(2025, 10, 3, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 19, 0, 0, 0, time.UTC)),
			description:   "月日曜日+時刻、スペースなし",
		},
		{
			no:            32,
			name:          "月日曜日+時刻（半角スペース）",
			content:       "10月3日（月） 18:00",
			expectedStart: time.Date(2025, 10, 3, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 19, 0, 0, 0, time.UTC)),
			description:   "月日曜日+時刻、半角スペース",
		},
		{
			no:            33,
			name:          "月日曜日+時刻（全角スペース）",
			content:       "10月3日（月）　18:00",
			expectedStart: time.Date(2025, 10, 3, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 19, 0, 0, 0, time.UTC)),
			description:   "月日曜日+時刻、全角スペース",
		},

		// 34-39: 月日形式（曜日なし）
		{
			no:            34,
			name:          "月日+時刻範囲（スペースなし）",
			content:       "10月3日14:00-16:00",
			expectedStart: time.Date(2025, 10, 3, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 16, 0, 0, 0, time.UTC)),
			description:   "月日+時刻範囲、スペースなし",
		},
		{
			no:            35,
			name:          "月日+時刻範囲（半角スペース）",
			content:       "10月3日 14:00-16:00",
			expectedStart: time.Date(2025, 10, 3, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 16, 0, 0, 0, time.UTC)),
			description:   "月日+時刻範囲、半角スペース",
		},
		{
			no:            36,
			name:          "月日+時刻範囲（全角スペース）",
			content:       "10月3日　14:00-16:00",
			expectedStart: time.Date(2025, 10, 3, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 16, 0, 0, 0, time.UTC)),
			description:   "月日+時刻範囲、全角スペース",
		},
		{
			no:            37,
			name:          "月日+時刻（スペースなし）",
			content:       "10月3日14:00",
			expectedStart: time.Date(2025, 10, 3, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 15, 0, 0, 0, time.UTC)),
			description:   "月日+時刻、スペースなし",
		},
		{
			no:            38,
			name:          "月日+時刻（半角スペース）",
			content:       "10月3日 14:00",
			expectedStart: time.Date(2025, 10, 3, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 15, 0, 0, 0, time.UTC)),
			description:   "月日+時刻、半角スペース",
		},
		{
			no:            39,
			name:          "月日+時刻（全角スペース）",
			content:       "10月3日　14:00",
			expectedStart: time.Date(2025, 10, 3, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 15, 0, 0, 0, time.UTC)),
			description:   "月日+時刻、全角スペース",
		},

		// 40-43: 時刻のみパターン
		{
			no:            40,
			name:          "時刻範囲（ハイフン）",
			content:       "14:00-16:00",
			expectedStart: time.Date(2025, 10, 3, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 16, 0, 0, 0, time.UTC)),
			description:   "時刻範囲、ハイフン区切り",
		},
		{
			no:            41,
			name:          "時刻範囲（半角波線）",
			content:       "14:00〜16:00",
			expectedStart: time.Date(2025, 10, 3, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 16, 0, 0, 0, time.UTC)),
			description:   "時刻範囲、半角波線区切り",
		},
		{
			no:            42,
			name:          "時刻範囲（全角波線）",
			content:       "14:00～16:00",
			expectedStart: time.Date(2025, 10, 3, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 16, 0, 0, 0, time.UTC)),
			description:   "時刻範囲、全角波線区切り",
		},
		{
			no:            43,
			name:          "時刻のみ",
			content:       "14:00",
			expectedStart: time.Date(2025, 10, 3, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 15, 0, 0, 0, time.UTC)),
			description:   "時刻のみ",
		},

		// 44-50: 日本語時分表現
		{
			no:            44,
			name:          "時分範囲（半角波線）",
			content:       "14時30分〜16時45分",
			expectedStart: time.Date(2025, 10, 3, 14, 30, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 16, 45, 0, 0, time.UTC)),
			description:   "日本語時分範囲、半角波線",
		},
		{
			no:            45,
			name:          "時分範囲（全角波線）",
			content:       "14時30分～16時45分",
			expectedStart: time.Date(2025, 10, 3, 14, 30, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 16, 45, 0, 0, time.UTC)),
			description:   "日本語時分範囲、全角波線",
		},
		{
			no:            46,
			name:          "時分範囲（から）",
			content:       "14時30分から16時45分",
			expectedStart: time.Date(2025, 10, 3, 14, 30, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 16, 45, 0, 0, time.UTC)),
			description:   "日本語時分範囲、から",
		},
		{
			no:            47,
			name:          "時分開始（から）",
			content:       "14時30分から",
			expectedStart: time.Date(2025, 10, 3, 14, 30, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 15, 30, 0, 0, time.UTC)),
			description:   "日本語時分開始、から",
		},
		{
			no:            48,
			name:          "時分開始（半角波線）",
			content:       "14時30分〜",
			expectedStart: time.Date(2025, 10, 3, 14, 30, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 15, 30, 0, 0, time.UTC)),
			description:   "日本語時分開始、半角波線",
		},
		{
			no:            49,
			name:          "時分開始（全角波線）",
			content:       "14時30分～",
			expectedStart: time.Date(2025, 10, 3, 14, 30, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 15, 30, 0, 0, time.UTC)),
			description:   "日本語時分開始、全角波線",
		},
		{
			no:            50,
			name:          "時分のみ",
			content:       "14時30分",
			expectedStart: time.Date(2025, 10, 3, 14, 30, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 15, 30, 0, 0, time.UTC)),
			description:   "日本語時分のみ",
		},

		// 51-56: 日本語時刻表現
		{
			no:            51,
			name:          "時刻範囲（半角波線）",
			content:       "14時〜16時",
			expectedStart: time.Date(2025, 10, 3, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 16, 0, 0, 0, time.UTC)),
			description:   "日本語時刻範囲、半角波線",
		},
		{
			no:            52,
			name:          "時刻範囲（全角波線）",
			content:       "14時～16時",
			expectedStart: time.Date(2025, 10, 3, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 16, 0, 0, 0, time.UTC)),
			description:   "日本語時刻範囲、全角波線",
		},
		{
			no:            53,
			name:          "時刻範囲（から）",
			content:       "14時から16時",
			expectedStart: time.Date(2025, 10, 3, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 16, 0, 0, 0, time.UTC)),
			description:   "日本語時刻範囲、から",
		},
		{
			no:            54,
			name:          "時刻開始（から）",
			content:       "14時から",
			expectedStart: time.Date(2025, 10, 3, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 15, 0, 0, 0, time.UTC)),
			description:   "日本語時刻開始、から",
		},
		{
			no:            55,
			name:          "時刻開始（半角波線）",
			content:       "14時〜",
			expectedStart: time.Date(2025, 10, 3, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 15, 0, 0, 0, time.UTC)),
			description:   "日本語時刻開始、半角波線",
		},
		{
			no:            56,
			name:          "時刻開始（全角波線）",
			content:       "14時～",
			expectedStart: time.Date(2025, 10, 3, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 15, 0, 0, 0, time.UTC)),
			description:   "日本語時刻開始、全角波線",
		},

		// 57-65: 相対日付表現
		{
			no:            57,
			name:          "明日+時刻（スペースなし）",
			content:       "明日14:00",
			expectedStart: time.Date(2025, 10, 4, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 4, 15, 0, 0, 0, time.UTC)),
			description:   "明日+時刻、スペースなし",
		},
		{
			no:            58,
			name:          "明日+時刻（半角スペース）",
			content:       "明日 14:00",
			expectedStart: time.Date(2025, 10, 4, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 4, 15, 0, 0, 0, time.UTC)),
			description:   "明日+時刻、半角スペース",
		},
		{
			no:            59,
			name:          "明日+時刻（全角スペース）",
			content:       "明日　14:00",
			expectedStart: time.Date(2025, 10, 4, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 4, 15, 0, 0, 0, time.UTC)),
			description:   "明日+時刻、全角スペース",
		},
		{
			no:            60,
			name:          "今日+時刻（スペースなし）",
			content:       "今日14:00",
			expectedStart: time.Date(2025, 10, 3, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 15, 0, 0, 0, time.UTC)),
			description:   "今日+時刻、スペースなし",
		},
		{
			no:            61,
			name:          "今日+時刻（半角スペース）",
			content:       "今日 18:00",
			expectedStart: time.Date(2025, 10, 3, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 19, 0, 0, 0, time.UTC)),
			description:   "今日+時刻、半角スペース",
		},
		{
			no:            62,
			name:          "今日+時刻（全角スペース）",
			content:       "今日　18:00",
			expectedStart: time.Date(2025, 10, 3, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 19, 0, 0, 0, time.UTC)),
			description:   "今日+時刻、全角スペース",
		},
		{
			no:            63,
			name:          "明後日+時刻（スペースなし）",
			content:       "明後日10:30",
			expectedStart: time.Date(2025, 10, 5, 10, 30, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 5, 11, 30, 0, 0, time.UTC)),
			description:   "明後日+時刻、スペースなし",
		},
		{
			no:            64,
			name:          "明後日+時刻（半角スペース）",
			content:       "明後日 10:30",
			expectedStart: time.Date(2025, 10, 5, 10, 30, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 5, 11, 30, 0, 0, time.UTC)),
			description:   "明後日+時刻、半角スペース",
		},
		{
			no:            65,
			name:          "明後日+時刻（全角スペース）",
			content:       "明後日　10:30",
			expectedStart: time.Date(2025, 10, 5, 10, 30, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 5, 11, 30, 0, 0, time.UTC)),
			description:   "明後日+時刻、全角スペース",
		},

		// 66-70: 時間帯表現
		{
			no:            66,
			name:          "午前",
			content:       "午前10時",
			expectedStart: time.Date(2025, 10, 3, 10, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 11, 0, 0, 0, time.UTC)),
			description:   "午前時刻",
		},
		{
			no:            67,
			name:          "午後",
			content:       "午後3時",
			expectedStart: time.Date(2025, 10, 3, 15, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 16, 0, 0, 0, time.UTC)),
			description:   "午後時刻",
		},
		{
			no:            68,
			name:          "夜",
			content:       "夜8時",
			expectedStart: time.Date(2025, 10, 3, 20, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 21, 0, 0, 0, time.UTC)),
			description:   "夜時刻",
		},
		{
			no:            69,
			name:          "朝",
			content:       "朝9時",
			expectedStart: time.Date(2025, 10, 3, 9, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 10, 0, 0, 0, time.UTC)),
			description:   "朝時刻",
		},
		{
			no:            70,
			name:          "昼",
			content:       "昼12時",
			expectedStart: time.Date(2025, 10, 3, 12, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 13, 0, 0, 0, time.UTC)),
			description:   "昼時刻",
		},

		// 71-76: 英語混在表現
		{
			no:            71,
			name:          "AM（スペースなし）",
			content:       "AM9:30",
			expectedStart: time.Date(2025, 10, 3, 9, 30, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 10, 30, 0, 0, time.UTC)),
			description:   "AM時刻、スペースなし",
		},
		{
			no:            72,
			name:          "AM（半角スペース）",
			content:       "AM 9:30",
			expectedStart: time.Date(2025, 10, 3, 9, 30, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 10, 30, 0, 0, time.UTC)),
			description:   "AM時刻、半角スペース",
		},
		{
			no:            73,
			name:          "AM（全角スペース）",
			content:       "AM　9:30",
			expectedStart: time.Date(2025, 10, 3, 9, 30, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 10, 30, 0, 0, time.UTC)),
			description:   "AM時刻、全角スペース",
		},
		{
			no:            74,
			name:          "PM（スペースなし）",
			content:       "PM6:15",
			expectedStart: time.Date(2025, 10, 3, 18, 15, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 19, 15, 0, 0, time.UTC)),
			description:   "PM時刻、スペースなし",
		},
		{
			no:            75,
			name:          "PM（半角スペース）",
			content:       "PM 6:15",
			expectedStart: time.Date(2025, 10, 3, 18, 15, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 19, 15, 0, 0, time.UTC)),
			description:   "PM時刻、半角スペース",
		},
		{
			no:            76,
			name:          "PM（全角スペース）",
			content:       "PM　6:15",
			expectedStart: time.Date(2025, 10, 3, 18, 15, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 19, 15, 0, 0, time.UTC)),
			description:   "PM時刻、全角スペース",
		},

		// 77-83: 区切り文字バリエーション
		{
			no:            77,
			name:          "ハイフン区切り（半角スペース）",
			content:       "10-6 18:00",
			expectedStart: time.Date(2025, 10, 6, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 19, 0, 0, 0, time.UTC)),
			description:   "ハイフン区切り日付、半角スペース",
		},
		{
			no:            78,
			name:          "ハイフン区切り（全角スペース）",
			content:       "10-6　18:00",
			expectedStart: time.Date(2025, 10, 6, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 19, 0, 0, 0, time.UTC)),
			description:   "ハイフン区切り日付、全角スペース",
		},
		{
			no:            79,
			name:          "ドット区切り（半角スペース）",
			content:       "10.6 18:00",
			expectedStart: time.Date(2025, 10, 6, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 19, 0, 0, 0, time.UTC)),
			description:   "ドット区切り日付、半角スペース",
		},
		{
			no:            80,
			name:          "ドット区切り（全角スペース）",
			content:       "10.6　18:00",
			expectedStart: time.Date(2025, 10, 6, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 19, 0, 0, 0, time.UTC)),
			description:   "ドット区切り日付、全角スペース",
		},
		{
			no:            81,
			name:          "完全スラッシュ（スペースなし）",
			content:       "2025/12/25 20:00",
			expectedStart: time.Date(2025, 12, 25, 20, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 12, 25, 21, 0, 0, 0, time.UTC)),
			description:   "完全スラッシュ日付、スペースなし",
		},
		{
			no:            82,
			name:          "完全スラッシュ（半角スペース）",
			content:       "2025/12/25 20:00",
			expectedStart: time.Date(2025, 12, 25, 20, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 12, 25, 21, 0, 0, 0, time.UTC)),
			description:   "完全スラッシュ日付、半角スペース",
		},
		{
			no:            83,
			name:          "完全スラッシュ（全角スペース）",
			content:       "2025/12/25　20:00",
			expectedStart: time.Date(2025, 12, 25, 20, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 12, 25, 21, 0, 0, 0, time.UTC)),
			description:   "完全スラッシュ日付、全角スペース",
		},

		// 84-89: 曖昧な時間表現
		{
			no:            84,
			name:          "夕方",
			content:       "夕方開始",
			expectedStart: time.Date(2025, 10, 3, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 20, 0, 0, 0, time.UTC)),
			description:   "夕方時間帯",
		},
		{
			no:            85,
			name:          "お昼頃",
			content:       "お昼頃開始",
			expectedStart: time.Date(2025, 10, 3, 12, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 13, 0, 0, 0, time.UTC)),
			description:   "お昼頃時間帯",
		},
		{
			no:            86,
			name:          "昼頃",
			content:       "昼頃開始",
			expectedStart: time.Date(2025, 10, 3, 12, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 13, 0, 0, 0, time.UTC)),
			description:   "昼頃時間帯",
		},
		{
			no:            87,
			name:          "深夜",
			content:       "深夜開始",
			expectedStart: time.Date(2025, 10, 3, 0, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 2, 0, 0, 0, time.UTC)),
			description:   "深夜時間帯",
		},
		{
			no:            88,
			name:          "夜中",
			content:       "夜中開始",
			expectedStart: time.Date(2025, 10, 3, 0, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 2, 0, 0, 0, time.UTC)),
			description:   "夜中時間帯",
		},
		{
			no:            89,
			name:          "早朝",
			content:       "早朝開始",
			expectedStart: time.Date(2025, 10, 3, 6, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 8, 0, 0, 0, time.UTC)),
			description:   "早朝時間帯",
		},

		// 90-103: 自然な日本語表現
		{
			no:            90,
			name:          "今週の曜日+時刻（スペースなし）",
			content:       "今週の土曜日14:00",
			expectedStart: time.Date(2025, 10, 4, 14, 0, 0, 0, time.UTC), // 明日の土曜日
			expectedEnd:   timePtr(time.Date(2025, 10, 4, 15, 0, 0, 0, time.UTC)),
			description:   "今週の曜日+時刻、スペースなし",
		},
		{
			no:            91,
			name:          "今週の曜日+時刻（半角スペース）",
			content:       "今週の土曜日 14:00",
			expectedStart: time.Date(2025, 10, 4, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 4, 15, 0, 0, 0, time.UTC)),
			description:   "今週の曜日+時刻、半角スペース",
		},
		{
			no:            92,
			name:          "今週の曜日+時刻（全角スペース）",
			content:       "今週の土曜日　14:00",
			expectedStart: time.Date(2025, 10, 4, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 4, 15, 0, 0, 0, time.UTC)),
			description:   "今週の曜日+時刻、全角スペース",
		},
		{
			no:            93,
			name:          "今週の曜日のみ（未来）",
			content:       "今週の土曜日",
			expectedStart: time.Date(2025, 10, 4, 0, 0, 0, 0, time.UTC), // 明日の土曜日
			expectedEnd:   timePtr(time.Date(2025, 10, 4, 1, 0, 0, 0, time.UTC)),
			description:   "今週の曜日のみ（未来の曜日）",
		},
		{
			no:            94,
			name:          "今週の曜日のみ（過去）",
			content:       "今週の月曜日",
			expectedStart: time.Date(2025, 9, 29, 0, 0, 0, 0, time.UTC), // 今週の月曜日（過去）
			expectedEnd:   timePtr(time.Date(2025, 9, 29, 1, 0, 0, 0, time.UTC)),
			description:   "今週の曜日のみ（過去の曜日）",
		},
		{
			no:            95,
			name:          "今度の曜日+時刻（スペースなし）",
			content:       "今度の土曜日14:00",
			expectedStart: time.Date(2025, 10, 4, 14, 0, 0, 0, time.UTC), // 次の土曜日は10/4
			expectedEnd:   timePtr(time.Date(2025, 10, 4, 15, 0, 0, 0, time.UTC)),
			description:   "今度の曜日+時刻、スペースなし",
		},
		{
			no:            96,
			name:          "今度の曜日+時刻（半角スペース）",
			content:       "今度の土曜日 14:00",
			expectedStart: time.Date(2025, 10, 4, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 4, 15, 0, 0, 0, time.UTC)),
			description:   "今度の曜日+時刻、半角スペース",
		},
		{
			no:            97,
			name:          "今度の曜日+時刻（全角スペース）",
			content:       "今度の土曜日　14:00",
			expectedStart: time.Date(2025, 10, 4, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 4, 15, 0, 0, 0, time.UTC)),
			description:   "今度の曜日+時刻、全角スペース",
		},
		{
			no:            98,
			name:          "来週の曜日+時刻（スペースなし）",
			content:       "来週の土曜日14:00",
			expectedStart: time.Date(2025, 10, 11, 14, 0, 0, 0, time.UTC), // 来週の土曜日は10/11
			expectedEnd:   timePtr(time.Date(2025, 10, 11, 15, 0, 0, 0, time.UTC)),
			description:   "来週の曜日+時刻、スペースなし",
		},
		{
			no:            99,
			name:          "来週の曜日+時刻（半角スペース）",
			content:       "来週の土曜日 14:00",
			expectedStart: time.Date(2025, 10, 11, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 11, 15, 0, 0, 0, time.UTC)),
			description:   "来週の曜日+時刻、半角スペース",
		},
		{
			no:            100,
			name:          "来週の曜日+時刻（全角スペース）",
			content:       "来週の土曜日　14:00",
			expectedStart: time.Date(2025, 10, 11, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 11, 15, 0, 0, 0, time.UTC)),
			description:   "来週の曜日+時刻、全角スペース",
		},
		{
			no:            101,
			name:          "来週曜日+時刻（スペースなし）",
			content:       "来週月曜日10:00",
			expectedStart: time.Date(2025, 10, 13, 10, 0, 0, 0, time.UTC), // 来週月曜日は10/13
			expectedEnd:   timePtr(time.Date(2025, 10, 13, 11, 0, 0, 0, time.UTC)),
			description:   "来週曜日+時刻、スペースなし",
		},
		{
			no:            102,
			name:          "来週曜日+時刻（半角スペース）",
			content:       "来週月曜日 10:00",
			expectedStart: time.Date(2025, 10, 13, 10, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 13, 11, 0, 0, 0, time.UTC)),
			description:   "来週曜日+時刻、半角スペース",
		},
		{
			no:            103,
			name:          "来週曜日+時刻（全角スペース）",
			content:       "来週月曜日　10:00",
			expectedStart: time.Date(2025, 10, 13, 10, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 13, 11, 0, 0, 0, time.UTC)),
			description:   "来週曜日+時刻、全角スペース",
		},
		{
			no:            104,
			name:          "来週の曜日のみ",
			content:       "来週の土曜日",
			expectedStart: time.Date(2025, 10, 11, 0, 0, 0, 0, time.UTC), // 来週の土曜日は10/11
			expectedEnd:   timePtr(time.Date(2025, 10, 11, 1, 0, 0, 0, time.UTC)),
			description:   "来週の曜日のみ",
		},
		{
			no:            105,
			name:          "今度の曜日のみ",
			content:       "今度の土曜日",
			expectedStart: time.Date(2025, 10, 4, 0, 0, 0, 0, time.UTC), // 次の土曜日は10/4
			expectedEnd:   timePtr(time.Date(2025, 10, 4, 1, 0, 0, 0, time.UTC)),
			description:   "今度の曜日のみ",
		},
		{
			no:            106,
			name:          "来週曜日のみ",
			content:       "来週月曜日",
			expectedStart: time.Date(2025, 10, 13, 0, 0, 0, 0, time.UTC), // 来週月曜日は10/13
			expectedEnd:   timePtr(time.Date(2025, 10, 13, 1, 0, 0, 0, time.UTC)),
			description:   "来週曜日のみ",
		},

		// 107-108: 期間表現
		{
			no:            107,
			name:          "時間期間",
			content:       "3時間のイベント",
			expectedStart: baseTime,
			expectedEnd:   timePtr(baseTime.Add(3 * time.Hour)),
			description:   "時間期間表現",
		},
		{
			no:            108,
			name:          "分期間",
			content:       "30分間のセッション",
			expectedStart: baseTime,
			expectedEnd:   timePtr(baseTime.Add(30 * time.Minute)),
			description:   "分期間表現",
		},

		// 109: デフォルト
		{
			no:            109,
			name:          "パターンなし",
			content:       "イベントです",
			expectedStart: time.Date(2025, 10, 3, 0, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 3, 1, 0, 0, 0, time.UTC)),
			description:   "パターンなし（デフォルト）",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualStart, actualEnd, _ := service.ExtractDateTime(tt.content, baseTime)

			// 開始時刻の検証
			if !actualStart.Equal(tt.expectedStart) {
				t.Errorf("No.%d %s - 開始時刻が一致しません\n期待値: %s\n実際値: %s\n説明: %s",
					tt.no, tt.name,
					tt.expectedStart.Format("2006-01-02 15:04:05"),
					actualStart.Format("2006-01-02 15:04:05"),
					tt.description)
			}

			// 終了時刻の検証
			if (tt.expectedEnd == nil && actualEnd != nil) ||
				(tt.expectedEnd != nil && actualEnd == nil) ||
				(tt.expectedEnd != nil && actualEnd != nil && !actualEnd.Equal(*tt.expectedEnd)) {

				expectedEndStr := "nil"
				if tt.expectedEnd != nil {
					expectedEndStr = tt.expectedEnd.Format("2006-01-02 15:04:05")
				}
				actualEndStr := "nil"
				if actualEnd != nil {
					actualEndStr = actualEnd.Format("2006-01-02 15:04:05")
				}

				t.Errorf("No.%d %s - 終了時刻が一致しません\n期待値: %s\n実際値: %s\n説明: %s",
					tt.no, tt.name,
					expectedEndStr, actualEndStr, tt.description)
			}

			// 成功ログ
			t.Logf("✅ No.%d %s: %s → %s to %s",
				tt.no, tt.name, tt.content,
				actualStart.Format("2006-01-02 15:04:05"),
				func() string {
					if actualEnd != nil {
						return actualEnd.Format("2006-01-02 15:04:05")
					}
					return "nil"
				}())
		})
	}
}

// パターン優先順位のテスト
func TestDateTimeExtractionService_PatternPriority(t *testing.T) {
	service := NewDateTimeExtractionService()
	baseTime := time.Date(2025, 10, 3, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		content       string
		expectedStart time.Time
		expectedEnd   *time.Time
		description   string
	}{
		{
			name:          "複数パターン（年月日優先）",
			content:       "2026年1月10日 14:00のイベントです。詳細は14:00から。",
			expectedStart: time.Date(2026, 1, 10, 14, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2026, 1, 10, 15, 0, 0, 0, time.UTC)),
			description:   "より具体的なパターンが優先されることを確認",
		},
		{
			name:          "曜日付きと曜日なし（曜日付き優先）",
			content:       "10/6（月）18:00に開催。別件で10/6 20:00。",
			expectedStart: time.Date(2025, 10, 6, 18, 0, 0, 0, time.UTC),
			expectedEnd:   timePtr(time.Date(2025, 10, 6, 19, 0, 0, 0, time.UTC)),
			description:   "曜日付きパターンが優先されることを確認",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualStart, actualEnd, _ := service.ExtractDateTime(tt.content, baseTime)

			if !actualStart.Equal(tt.expectedStart) {
				t.Errorf("開始時刻が一致しません\n期待値: %s\n実際値: %s\n説明: %s",
					tt.expectedStart.Format("2006-01-02 15:04:05"),
					actualStart.Format("2006-01-02 15:04:05"),
					tt.description)
			}

			if (tt.expectedEnd == nil && actualEnd != nil) ||
				(tt.expectedEnd != nil && actualEnd == nil) ||
				(tt.expectedEnd != nil && actualEnd != nil && !actualEnd.Equal(*tt.expectedEnd)) {

				expectedEndStr := "nil"
				if tt.expectedEnd != nil {
					expectedEndStr = tt.expectedEnd.Format("2006-01-02 15:04:05")
				}
				actualEndStr := "nil"
				if actualEnd != nil {
					actualEndStr = actualEnd.Format("2006-01-02 15:04:05")
				}

				t.Errorf("終了時刻が一致しません\n期待値: %s\n実際値: %s\n説明: %s",
					expectedEndStr, actualEndStr, tt.description)
			}
		})
	}
}

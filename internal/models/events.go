package models

import "time"

// 各推しのイベント情報
type Event struct {
	ID                    int64         `json:"id"`
	Title                 string        `json:"title"`
	Description           *string       `json:"description"`
	URL                   *string       `json:"url"`
	Starts_at             time.Time     `json:"starts_at"`
	Ends_at               *time.Time    `json:"ends_at"`
	Has_alarm             bool          `json:"has_alarm"`
	Notification_timing   string        `json:"notification_timing"`
	Has_notification_sent bool          `json:"has_notification_sent"`
	Category              *CategoryItem `json:"category"`
}

// 各推しのイベントのカテゴリ情報
type CategoryItem struct {
	ID   int64  `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// 推しのイベント情報レスポンス
type OshiEventsResponse struct {
	Oshis []OshiEventsResponseItem `json:"oshis"`
}

// 推しのイベント情報レスポンス内の各推し情報
type OshiEventsResponseItem struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	Color  string  `json:"color"`
	Events []Event `json:"events"`
}

// イベント詳細情報
type EventDetail struct {
	ID                    int64            `json:"id"`
	Title                 string           `json:"title"`
	Description           *string          `json:"description"`
	URL                   *string          `json:"url"`
	Starts_at             time.Time        `json:"starts_at"`
	Ends_at               *time.Time       `json:"ends_at"`
	Has_alarm             bool             `json:"has_alarm"`
	Notification_timing   string           `json:"notification_timing"`
	Has_notification_sent bool             `json:"has_notification_sent"`
	Oshi                  EventOshi        `json:"oshi"`
}

// イベント詳細レスポンス用の推し情報
type EventOshi struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// イベント詳細レスポンス
type EventDetailResponse struct {
	Event EventDetail `json:"event"`
}

// イベント更新リクエスト
type UpdateEventRequest struct {
	Event UpdateEventData `json:"event"`
}

// イベント更新データ
type UpdateEventData struct {
	Title               string     `json:"title" validate:"required"`
	Description         *string    `json:"description"`
	URL                 *string    `json:"url"`
	Starts_at           time.Time  `json:"starts_at" validate:"required"`
	Ends_at             *time.Time `json:"ends_at"`
	Has_alarm           bool       `json:"has_alarm"`
	Notification_timing string     `json:"notification_timing" validate:"required"`
}

// イベント更新レスポンス
type UpdateEventResponse struct {
	Event UpdatedEventDetail `json:"event"`
}

// 更新されたイベント詳細情報（推し情報なし）
type UpdatedEventDetail struct {
	ID                    int64      `json:"id"`
	Title                 string     `json:"title"`
	Description           *string    `json:"description"`
	URL                   *string    `json:"url"`
	Starts_at             time.Time  `json:"starts_at"`
	Ends_at               *time.Time `json:"ends_at"`
	Has_alarm             bool       `json:"has_alarm"`
	Notification_timing   string     `json:"notification_timing"`
	Has_notification_sent bool       `json:"has_notification_sent"`
}

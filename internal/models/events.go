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

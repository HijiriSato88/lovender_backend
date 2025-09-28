package models

import "time"

// カテゴリ情報
type Category struct {
	ID          uint16    `json:"id" db:"id"`
	Slug        string    `json:"slug" db:"slug"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	CreatedAt   time.Time `json:"-" db:"created_at"`
	UpdatedAt   time.Time `json:"-" db:"updated_at"`
}

// カテゴリ一覧レスポンス
type CommonResponse struct {
	Categories []Category `json:"categories"`
}

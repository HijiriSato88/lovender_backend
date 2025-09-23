package models

import "time"

// カテゴリ情報
type Category struct {
	ID          *int16    `json:"id" db:"id"`
	Slug        *string   `json:"slug" db:"slug"`
	Name        *string   `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

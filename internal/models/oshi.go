package models

import "time"

// 推し情報
type Oshi struct {
	ID          int64     `json:"id" db:"id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	ThemeColor  string    `json:"color" db:"theme_color"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// 推しのアカウント情報
type OshiAccount struct {
	ID        *int64    `json:"id" db:"id"`
	OshiID    int64     `json:"oshi_id" db:"oshi_id"`
	URL       *string   `json:"url" db:"url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// 推しレスポンス
type OshiResponse struct {
	ID         int64    `json:"id"`
	Name       string   `json:"name"`
	Color      string   `json:"color"`
	URLs       []string `json:"urls"`
	Categories []string `json:"categories"`
}

// 推し詳細情報
type OshiWithDetails struct {
	Oshi       *Oshi
	Accounts   []*OshiAccount
	Categories []*Category
}

// 推し一覧レスポンス
type OshisResponse struct {
	Oshis []OshiResponse `json:"oshis"`
}

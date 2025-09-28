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
	ID        int64     `json:"id" db:"id"`
	OshiID    int64     `json:"oshi_id" db:"oshi_id"`
	URL       string    `json:"url" db:"url"`
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

// 推し作成リクエスト
type CreateOshiRequest struct {
	Name       string   `json:"name" validate:"required"`
	Color      string   `json:"color" validate:"required,startswith=#,len=7"`
	URLs       []string `json:"urls"`
	Categories []string `json:"categories"`
}

// 推し作成レスポンス
type CreateOshiResponse struct {
	Oshi CreateOshiResponseItem `json:"oshi"`
}

// 推し作成レスポンス内の推し情報
type CreateOshiResponseItem struct {
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

// 推し更新リクエスト
type UpdateOshiRequest struct {
	Name       string   `json:"name" validate:"required"`
	Color      string   `json:"color" validate:"required,startswith=#,len=7"`
	URLs       []string `json:"urls"`
	Categories []string `json:"categories"`
}

// 推し更新レスポンス
type UpdateOshiResponse struct {
	Oshi UpdateOshiResponseItem `json:"oshi"`
}

// 推し更新レスポンス内の推し情報
type UpdateOshiResponseItem struct {
	ID         int64    `json:"id"`
	Name       string   `json:"name"`
	Color      string   `json:"color"`
	URLs       []string `json:"urls"`
	Categories []string `json:"categories"`
}

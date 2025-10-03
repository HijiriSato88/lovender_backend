package models

import "time"

// 外部投稿API用のレスポンス構造体
type ExternalPostsResponse struct {
	Posts []ExternalPost `json:"posts"`
}

// 外部投稿情報
type ExternalPost struct {
	ID        int64            `json:"id"`
	UserID    int64            `json:"userId"`
	Content   string           `json:"content"`
	CreatedAt string           `json:"createdAt"` // "2025-09-28 04:39:15" 形式
	User      ExternalPostUser `json:"user"`
}

// 外部投稿のユーザー情報
type ExternalPostUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatarUrl"`
}

// 投稿からイベント作成用の内部構造体
type PostEventCandidate struct {
	OshiID      int64
	PostID      int64
	Content     string
	CreatedAt   time.Time
	AccountName string
	Keywords    []string // マッチしたキーワード
}

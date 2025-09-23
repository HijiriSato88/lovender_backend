package models

import "time"

// 新規登録リクエスト
type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// 新規登録レスポンス
type RegisterResponse struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Token     string    `json:"token"`
}

// ログインリクエスト
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// ログインレスポンス
type LoginResponse struct {
	Token string `json:"token"`
}

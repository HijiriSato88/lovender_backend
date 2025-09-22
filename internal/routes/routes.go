package routes

import (
	"lovender_backend/internal/handler"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, userHandler *handler.UserHandler) {
	// API ルート
	api := e.Group("/api")

	// ユーザー情報取得
	api.GET("/users/:id", userHandler.GetUser)
}

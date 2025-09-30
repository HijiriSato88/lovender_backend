package routes

import (
	"lovender_backend/internal/handler"
	"lovender_backend/pkg/jwtutil"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, userHandler *handler.UserHandler, oshiHandler *handler.OshiHandler, commonHandler *handler.CommonHandler, eventsHandler *handler.EventsHandler) {
	// ルート
	api := e.Group("/api")

	// 認証
	api.POST("/auth/register", userHandler.Register)
	api.POST("/auth/login", userHandler.Login)

	// 共通情報
	api.GET("/common", commonHandler.GetCommon)

	// JWT認証が必要なエンドポイント
	protected := api.Group("/me")
	protected.Use(jwtutil.JWTMiddleware())

	// 推し関連のエンドポイント
	protected.GET("/oshis", oshiHandler.GetMyOshis)
	protected.POST("/oshis/new", oshiHandler.CreateOshi)
	protected.PUT("/oshis/:oshiId", oshiHandler.UpdateOshi)

	// イベント関連のエンドポイント
	protected.GET("/events", eventsHandler.GetMyOshiEvents)
	protected.GET("/events/:eventId", eventsHandler.GetEventByID)
	protected.PUT("/events/:eventId", eventsHandler.UpdateEvent)
	protected.POST("/events/new", eventsHandler.CreateEvent)

	// API接続テスト用のユーザー情報取得
	api.GET("/users/:id", userHandler.GetUser)
}

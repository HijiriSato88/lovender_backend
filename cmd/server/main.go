package main

import (
	"context"
	"lovender_backend/internal/cache"
	"lovender_backend/internal/database"
	"lovender_backend/internal/handler"
	"lovender_backend/internal/repository"
	"lovender_backend/internal/routes"
	"lovender_backend/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// データベース接続
	db, err := database.NewConnection()
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	defer db.Close()

	// キャッシュマネージャーを初期化
	// 起動時にキーワードをメモリにロード
	cacheManager := cache.NewCacheManager(db)

	// 依存関係の注入
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	oshiRepo := repository.NewOshiRepository(db)
	oshiService := service.NewOshiService(oshiRepo)
	oshiHandler := handler.NewOshiHandler(oshiService)

	categoryRepo := repository.NewCategoryRepository(db)
	commonService := service.NewCommonService(categoryRepo)
	commonHandler := handler.NewCommonHandler(commonService)

	eventsRepo := repository.NewEventsRepository(db)
	eventsService := service.NewEventsService(eventsRepo)
	eventsHandler := handler.NewEventsHandler(eventsService)

	// イベント自動登録サービス
	eventAutoService := service.NewEventAutoService(eventsRepo, cacheManager.GetKeywordCache())
	eventAutoHandler := handler.NewEventAutoHandler(eventAutoService)

	// Echo インスタンスを作成
	e := echo.New()

	// ミドルウェア
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// ルート設定
	routes.SetupRoutes(e, userHandler, oshiHandler, commonHandler, eventsHandler, eventAutoHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Graceful shutdown
	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	e.Logger.Info("Server is shutting down...")

	// キャッシュマネージャーのシャットダウン
	cacheManager.Shutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	e.Logger.Info("Server stopped")
}

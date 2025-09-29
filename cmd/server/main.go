package main

import (
	"lovender_backend/internal/database"
	"lovender_backend/internal/handler"
	"lovender_backend/internal/repository"
	"lovender_backend/internal/routes"
	"lovender_backend/internal/service"
	"os"

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

	// Echo インスタンスを作成
	e := echo.New()

	// ミドルウェア
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// ルート設定
	routes.SetupRoutes(e, userHandler, oshiHandler, commonHandler, eventsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// サーバー起動
	e.Logger.Fatal(e.Start(":" + port))
}

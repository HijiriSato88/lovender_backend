package handler

import (
	"context"
	"lovender_backend/internal/service"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// EventAutoHandler イベント自動登録ハンドラー
type EventAutoHandler struct {
	eventAutoService *service.EventAutoService
}

// NewEventAutoHandler コンストラクタ
func NewEventAutoHandler(eventAutoService *service.EventAutoService) *EventAutoHandler {
	return &EventAutoHandler{
		eventAutoService: eventAutoService,
	}
}

// 自動イベント作成処理
func (h *EventAutoHandler) ProcessAutoEvents(c echo.Context) error {
	// タイムアウト付きのコンテキスト（5分）
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Minute)
	defer cancel()

	// 自動イベント作成処理を実行
	result, err := h.eventAutoService.ProcessAutoEventCreation(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to process auto events",
			"details": err.Error(),
		})
	}

	// 成功レスポンス
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Auto event processing completed",
		"result":  result,
	})
}

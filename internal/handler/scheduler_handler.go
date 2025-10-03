package handler

import (
	"lovender_backend/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

// スケジューラーハンドラー
type SchedulerHandler struct {
	schedulerService *service.SchedulerService
}

// コンストラクタ
func NewSchedulerHandler(schedulerService *service.SchedulerService) *SchedulerHandler {
	return &SchedulerHandler{
		schedulerService: schedulerService,
	}
}

// スケジューラーの状態を取得
func (h *SchedulerHandler) GetSchedulerStatus(c echo.Context) error {
	status := h.schedulerService.GetStatus()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Scheduler status retrieved successfully",
		"status":  status,
	})
}

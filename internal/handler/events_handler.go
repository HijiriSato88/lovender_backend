package handler

import (
	"lovender_backend/internal/service"
	"lovender_backend/pkg/jwtutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

type EventsHandler struct {
	eventsService service.EventsService
}

func NewEventsHandler(eventsService service.EventsService) *EventsHandler {
	return &EventsHandler{
		eventsService: eventsService,
	}
}

func (h EventsHandler) GetMyOshiEvents(c echo.Context) error {
	// JWTトークンからユーザー情報を取得
	claims, err := jwtutil.ExtractUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
	}

	userID := int64(claims.UserID)

	// ユーザーの登録した各推しのイベントを全て取得
	events, err := h.eventsService.GetUserOshiEvents(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, events)
}

package handler

import (
	"lovender_backend/internal/service"
	"lovender_backend/pkg/jwtutil"
	"net/http"
	"strconv"

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

// 特定のイベント詳細を取得
func (h EventsHandler) GetEventByID(c echo.Context) error {
	// JWTトークンからユーザー情報を取得
	claims, err := jwtutil.ExtractUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
	}

	// パスパラメータからeventIdを取得
	eventIDStr := c.Param("eventId")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid event ID"})
	}

	userID := int64(claims.UserID)

	// イベント詳細を取得
	event, err := h.eventsService.GetEventByID(eventID, userID)
	if err != nil {
		if err.Error() == "event not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Event not found"})
		}
		if err.Error() == "access denied" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, event)
}

package handler

import (
	"lovender_backend/internal/models"
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

// イベント情報を更新
func (h EventsHandler) UpdateEvent(c echo.Context) error {
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

	// リクエストBodyのバインド
	var req models.UpdateEventRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// バリデーション
	if req.Event.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title is required"})
	}
	if req.Event.Notification_timing == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Notification timing is required"})
	}

	// 日時をUTCに変換
	req.Event.Starts_at = req.Event.Starts_at.UTC()
	if req.Event.Ends_at != nil {
		utcTime := req.Event.Ends_at.UTC()
		req.Event.Ends_at = &utcTime
	}

	userID := int64(claims.UserID)

	// イベント更新
	event, err := h.eventsService.UpdateEvent(eventID, userID, &req.Event)
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

// イベントを新規作成
func (h EventsHandler) CreateEvent(c echo.Context) error {
	// JWTトークンからユーザー情報を取得
	claims, err := jwtutil.ExtractUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
	}

	// リクエストBodyのバインド
	var req models.CreateEventRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// バリデーション
	if req.Event.OshiID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Oshi ID is required"})
	}
	if req.Event.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title is required"})
	}
	if req.Event.Notification_timing == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Notification timing is required"})
	}

	// 日時をUTCに変換
	req.Event.Starts_at = req.Event.Starts_at.UTC()
	if req.Event.Ends_at != nil {
		utcTime := req.Event.Ends_at.UTC()
		req.Event.Ends_at = &utcTime
	}

	userID := int64(claims.UserID)

	// イベント作成
	event, err := h.eventsService.CreateEvent(userID, &req.Event)
	if err != nil {
		if err.Error() == "oshi not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Oshi not found"})
		}
		if err.Error() == "access denied" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusCreated, event)
}

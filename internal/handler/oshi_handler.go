package handler

import (
	"lovender_backend/internal/service"
	"lovender_backend/pkg/jwtutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

type OshiHandler struct {
	oshiService service.OshiService
}

func NewOshiHandler(oshiService service.OshiService) *OshiHandler {
	return &OshiHandler{
		oshiService: oshiService,
	}
}

func (h *OshiHandler) GetMyOshis(c echo.Context) error {
	// JWTトークンからユーザー情報を取得
	claims, err := jwtutil.ExtractUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
	}

	userID := int64(claims.UserID)

	// ユーザーの推し一覧を取得
	oshis, err := h.oshiService.GetUserOshis(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, oshis)
}

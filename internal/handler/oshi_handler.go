package handler

import (
	"log"
	"lovender_backend/internal/models"
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

// 推しの新規作成
func (h *OshiHandler) CreateOshi(c echo.Context) error {
	// ユーザー情報を取得
	claims, err := jwtutil.ExtractUser(c)
	if err != nil {
		log.Printf("CreateOshi ERROR: invalid token: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
	}

	// リクエストBodyのバインド
	var req models.CreateOshiRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("CreateOshi ERROR: bind failed: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	resp, err := h.oshiService.CreateOshi(int64(claims.UserID), &req)
	if err != nil {
		if err.Error() == "oshi already exists" {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Oshi already exists"})
		}
		if err.Error() == "invalid categories provided" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create oshi"})
	}

	return c.JSON(http.StatusCreated, resp)
}

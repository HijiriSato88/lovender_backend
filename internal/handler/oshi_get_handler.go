package handler

import (
	"log"
	"lovender_backend/internal/service"
	"lovender_backend/pkg/jwtutil"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type OshiGetHandler struct {
	oshiGetService service.OshiGetService
}

func NewOshiGetHandler(oshiGetService service.OshiGetService) *OshiGetHandler {
	return &OshiGetHandler{
		oshiGetService: oshiGetService,
	}
}

func (h *OshiGetHandler) GetMyOshiByID(c echo.Context) error {
	//JWTからユーザーIDを取得
	claims, err := jwtutil.ExtractUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
	}

	userID := int64(claims.UserID)

	//パスパラメータからoshiIDを取得
	oshiIDStr := c.Param("oshiId")
	log.Printf("DEBUG: received oshiId param=%s", oshiIDStr)

	oshiID, err := strconv.ParseInt(oshiIDStr, 10, 64)
	if err != nil {
		log.Printf("GetOshi ERROR: invalid oshiId: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid oshi ID"})
	}

	//Serviceを呼び出して推し1人を取得
	resp, err := h.oshiGetService.GetOshiByID(oshiID, userID)
	if err != nil {
		if err.Error() == "oshi not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Oshi not found"})
		}
		log.Printf("GetOshi ERROR: service error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusOK, resp)
}

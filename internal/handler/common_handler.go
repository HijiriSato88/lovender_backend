package handler

import (
	"lovender_backend/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CommonHandler struct {
	svc service.CommonService
}

func NewCommonHandler(s service.CommonService) *CommonHandler {
	return &CommonHandler{svc: s}
}

func (h *CommonHandler) GetCommon(c echo.Context) error {
	resp, err := h.svc.GetCommon(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusOK, resp)
}

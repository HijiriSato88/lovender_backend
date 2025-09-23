package handler

import (
	"lovender_backend/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CommonHandler struct {
	commonService service.CommonService
}

func NewCommonHandler(commonService service.CommonService) *CommonHandler {
	return &CommonHandler{commonService: commonService}
}

func (h *CommonHandler) GetCommon(c echo.Context) error {
	resp, err := h.commonService.GetCommon()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusOK, resp)
}

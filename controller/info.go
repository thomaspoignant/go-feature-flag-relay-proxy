package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/thomaspoignant/go-feature-flag-relay-proxy/service"
	"net/http"
)

type info struct {
	monitoringService service.Monitoring
}

func NewInfo(monitoring service.Monitoring) Controller {
	return &info{
		monitoringService: monitoring,
	}
}

func (h *info) Handler(c echo.Context) error {
	return c.JSON(http.StatusOK, h.monitoringService.Info())
}

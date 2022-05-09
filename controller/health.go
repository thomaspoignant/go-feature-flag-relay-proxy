package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/thomaspoignant/go-feature-flag-relay-proxy/service"
	"net/http"
)

type health struct {
	monitoringService service.Monitoring
}

// NewHealth is a constructor to create a new controller for the health method
func NewHealth(monitoring service.Monitoring) Controller {
	return &health{
		monitoringService: monitoring,
	}
}

// Handler is the entry point for this API
func (h *health) Handler(c echo.Context) error {
	return c.JSON(http.StatusOK, h.monitoringService.Health())
}

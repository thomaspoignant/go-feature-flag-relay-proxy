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

// Handler is the entry point for the Info API
// @Summary      Info, give information about the instance of go-feature-flag relay proxy
// @Description  Info, give information about the instance of go-feature-flag relay proxy
// @Tags         monitoring
// @Produce      json
// @Success      200  {object}   model.InfoResponse
// @Router       /info [get]
func (h *info) Handler(c echo.Context) error {
	return c.JSON(http.StatusOK, h.monitoringService.Info())
}

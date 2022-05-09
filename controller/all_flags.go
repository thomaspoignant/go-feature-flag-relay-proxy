package controller

import (
	"github.com/labstack/echo/v4"
	ffclient "github.com/thomaspoignant/go-feature-flag"
)

type allFlags struct {
	goFF *ffclient.GoFeatureFlag
}

func NewAllFlags(goFF *ffclient.GoFeatureFlag) Controller {
	return &allFlags{
		goFF: goFF,
	}
}

func (h *allFlags) Handler(c echo.Context) error {
	// TODO:
	// - retriever the user
	// - validate the key
	// - call goFF
	// - return the results

	// return c.JSON(http.StatusOK, h.monitoringService.Health())
	return nil
}

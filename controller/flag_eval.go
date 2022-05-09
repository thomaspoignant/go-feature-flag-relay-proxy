package controller

import (
	"github.com/labstack/echo/v4"
	ffclient "github.com/thomaspoignant/go-feature-flag"
)

type flagEval struct {
	goFF *ffclient.GoFeatureFlag
}

func NewFlagEval(goFF *ffclient.GoFeatureFlag) Controller {
	return &flagEval{
		goFF: goFF,
	}
}

func (h *flagEval) Handler(c echo.Context) error {
	// TODO:
	// - retriever the user
	// - validate the key
	// - retrieve the flagKey
	// - call goFF
	// - return the results

	// return c.JSON(http.StatusOK, h.monitoringService.Health())
	return nil
}

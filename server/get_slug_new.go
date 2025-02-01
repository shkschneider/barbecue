package server

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"barbecue/api"
	"barbecue/core"
	"barbecue/driver"
)

func GetSlugNew(ctx echo.Context) error {
	out := driver.NewHtmlDriver(ctx)
	slug := ctx.Param("slug")
	parents, err := api.GetBySlug(slug)
	if err != nil {
		core.Log.Debug("GetSlugNew", err)
		return out.Err(http.StatusNotFound, "NotFound")
	}
	parent := (*parents)[0]
	return out.Out(driver.NewHtmlDriverData(http.StatusOK, driver.T_FORM, ApiResponse {
		Parent: &parent,
	}))
}

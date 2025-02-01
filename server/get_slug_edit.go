package server

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"barbecue/api"
	"barbecue/core"
	"barbecue/driver"
)

func GetSlugEdit(ctx echo.Context) error {
	out := driver.NewHtmlDriver(ctx)
	slug := ctx.Param("slug")
	tasks, err := api.GetBySlug(slug)
	if err != nil {
		core.Log.Debug("GetSlugEdit", err)
		return out.Err(http.StatusNotFound, "NotFound")
	}
	task := (*tasks)[0]
	return out.Out(driver.NewHtmlDriverData(http.StatusFound, driver.T_FORM, ApiResponse {
		Task: &task,
	}))
}

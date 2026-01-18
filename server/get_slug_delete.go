package server

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"barbecue/api"
	"barbecue/core"
	"barbecue/driver"
)

func GetSlugDelete(ctx echo.Context) error {
	out := driver.NewHtmlDriver(ctx)
	slug := ctx.Param("slug")
	tasks, err := api.GetBySlug(slug)
	if err != nil {
		core.Log.Debug("GetSlugDelete", err)
		return out.Err(http.StatusNotFound, "NotFound")
	}
	task := (*tasks)[0]
	if err := api.Remove(task) ; err != nil {
		core.Log.Debug("GetSlugDelete", err)
		return err
	}
	parent, err := api.GetParent(task)
	if err != nil {
		core.Log.Debug("GetSlugDelete", err)
	}
	if parent != nil {
		return out.Out(driver.NewHtmlDriverRedirect("/" + parent.Slug))
	}
	return out.Out(driver.NewHtmlDriverRedirect("/"))
}

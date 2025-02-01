package server

import (
	"net/http"
	"strconv"
	"github.com/labstack/echo/v4"
	"barbecue/api"
	"barbecue/core"
	"barbecue/driver"
)

func GetSlugProgress(ctx echo.Context) error {
	out := driver.NewHtmlDriver(ctx)
	slug := ctx.Param("slug")
	tasks, err := api.GetBySlug(slug)
	if err != nil {
		core.Log.Debug("GetSlugProgress", err)
		return out.Err(http.StatusNotFound, "NotFound")
	}
	task := (*tasks)[0]
	if pc, err := strconv.ParseUint(ctx.Param("progress"), 10, 16) ; err == nil {
		if pc <= 0 {
			task.Progress = 0
		} else if pc >= 100 {
			task.Progress = 100
		} else {
			task.Progress = uint(pc)
		}
	}
	if err := core.Database.Update(task) ; err != nil {
		core.Log.Debug("GetSlugProgress", err)
		return out.Err(http.StatusInternalServerError, "Database")
	}
	parent, err := api.GetParent(task)
	if err != nil {
		core.Log.Debug("GetSlugProgress", err)
		return out.Out(driver.NewHtmlDriverRedirect("/" + task.Slug))
	}
	return out.Out(driver.NewHtmlDriverRedirect("/" + parent.Slug))
}

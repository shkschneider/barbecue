package server

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"barbecue/api"
	"barbecue/core"
	"barbecue/driver"
)

func PostSlugEdit(ctx echo.Context) error {
	out := driver.NewHtmlDriver(ctx)
	slug := ctx.Param("slug")
	title := ctx.FormValue("title")
	description := ctx.FormValue("description")
	tasks, err := api.GetBySlug(slug)
	if err != nil {
		core.Log.Debug("PostSlugEdit", err)
		return out.Err(http.StatusNotFound, "NotFound")
	}
	task := (*tasks)[0]
	task.Title = title
	task.Description = description
	if _, err := api.Update(task) ; err != nil {
		core.Log.Debug("PostSlugEdit", err)
		return out.Err(http.StatusInternalServerError, "Database")
	}
	return out.Out(driver.NewHtmlDriverRedirect("/" + task.Slug))
}

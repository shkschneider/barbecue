package server

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"barbecue/api"
	"barbecue/core"
	"barbecue/driver"
)

func PostSlugNew(ctx echo.Context) error {
	out := driver.NewHtmlDriver(ctx)
	slug := ctx.Param("slug")
	title := ctx.FormValue("title")
	description := ctx.FormValue("description")
	parents, err := api.GetBySlug(slug)
	if err != nil {
		core.Log.Debug("PostSlugNew", err)
		return out.Err(http.StatusNotFound, "NotFound")
	}
	parent := (*parents)[0]
	tasks, err := api.Add(title, description)
	if err != nil {
		core.Log.Debug("PostSlugNew", err)
		return out.Err(http.StatusInternalServerError, "Database")
	}
	task := (*tasks)[0]
	task.Super = &parent.ID
	if err := core.Database.Update(task) ; err != nil {
		core.Log.Debug("PostSlugNew", err)
		return out.Err(http.StatusInternalServerError, "Database")
	}
	return out.Out(driver.NewHtmlDriverRedirect("/" + task.Slug))
}

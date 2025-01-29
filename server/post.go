package server

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"barbecue/api"
	"barbecue/core"
	"barbecue/driver"
)

func PostNew(ctx echo.Context) error {
	out := driver.NewHtmlDriver()
	tasks, err := api.Add(ctx.FormValue("title"), ctx.FormValue("description"))
	if err != nil { return out.Ko(ctx, http.StatusInternalServerError) }
	task := (*tasks)[0]
	core.Context.Logger.Debug(task)
	return out.Redirect(ctx, "/" + task.Slug)
}

func PostSlugNew(ctx echo.Context) error {
	out := driver.NewHtmlDriver()
	parents, err := api.GetBySlug(ctx.Param("slug"))
	if err != nil { return out.Ko(ctx, http.StatusNotFound) }
	parent := (*parents)[0]
	core.Context.Logger.Debug(parent)
	tasks, err := api.Add(ctx.FormValue("title"), ctx.FormValue("description"))
	if err != nil { return out.Ko(ctx, http.StatusInternalServerError) }
	task := (*tasks)[0]
	task.Super = &parent.ID
	core.Context.Logger.Debug(task)
	if err := core.Context.Database.Update(task) ; err != nil {
		core.Context.Logger.Error(err)
		return out.Ko(ctx, http.StatusInternalServerError)
	}
	return out.Redirect(ctx, "/" + task.Slug)
}

func PostSlugEdit(ctx echo.Context) error {
	out := driver.NewHtmlDriver()
	tasks, err := api.GetBySlug(ctx.Param("slug"))
	if err != nil { return out.Ko(ctx, http.StatusNotFound) }
	task := (*tasks)[0]
	task.Title = ctx.FormValue("title")
	task.Description = ctx.FormValue("description")
	core.Context.Logger.Debug(task)
	if _, err := api.Update(task) ; err != nil {
		core.Context.Logger.Error(err)
		return out.Ko(ctx, http.StatusInternalServerError)
	}
	return out.Redirect(ctx, "/" + task.Slug)
}

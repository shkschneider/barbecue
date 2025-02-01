package server

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"barbecue/api"
	"barbecue/core"
	"barbecue/driver"
)

func PostNew(ctx echo.Context) error {
	out := driver.NewHtmlDriver(ctx)
	title := ctx.FormValue("title")
	description := ctx.FormValue("description")
	tasks, err := api.Add(title, description)
	if err != nil {
		core.Log.Debug("PostNew", err)
		return out.Err(http.StatusInternalServerError, "Database")
	}
	task := (*tasks)[0]
	return out.Out(driver.NewHtmlDriverRedirect("/" + task.Slug))
}

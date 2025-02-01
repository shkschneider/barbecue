package server

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"barbecue/api"
	"barbecue/core"
	"barbecue/driver"
)

func GetSlug(ctx echo.Context) error {
	out := driver.NewHtmlDriver(ctx)
	tasks, err := api.GetBySlug(ctx.Param("slug"))
	if err != nil {
		core.Log.Debug("GetSlug", err)
		return out.Err(http.StatusNotFound, "NotFound")
	}
	task := (*tasks)[0]
	parent, _ := api.GetParent(task)
	children, _ := api.GetChildren(task)
	return out.Out(driver.NewHtmlDriverData(http.StatusFound, driver.T_TASK, ApiResponse {
		Parent: parent,
		Task: &task,
		Children: children,
	}))
}

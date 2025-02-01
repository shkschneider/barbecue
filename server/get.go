package server

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"barbecue/api"
	"barbecue/core"
	"barbecue/driver"
)

func Get(ctx echo.Context) error {
	out := driver.NewHtmlDriver(ctx)
	parents, err := api.GetParents()
	if err != nil {
		core.Log.Debug("Get", err)
		return out.Err(http.StatusNotFound, "NotFound")
	} else if parents == nil {
		core.Log.Debug("Get", core.ErrNothing)
		return out.Out(driver.NewHtmlDriverData(http.StatusOK, driver.T_INDEX, []ApiResponse{}))
	}
	var responses []ApiResponse
	for _, parent := range *parents {
		children, _ := api.GetChildren(parent)
		responses = append(responses, ApiResponse {
			Parent: nil,
			Task: &parent,
			Children: children,
		})
	}
	return out.Out(driver.NewHtmlDriverData(http.StatusOK, driver.T_INDEX, &responses))
}

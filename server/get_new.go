package server

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"barbecue/driver"
)

func GetNew(ctx echo.Context) error {
	out := driver.NewHtmlDriver(ctx)
	return out.Out(driver.NewHtmlDriverData(http.StatusOK, driver.T_FORM, nil))
}

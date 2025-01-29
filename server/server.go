package server

import (
	"fmt"
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/urfave/cli/v2"
	// "barbecue/api"
	"barbecue/core"
	"barbecue/driver"
)

func Serve(cli *cli.Context) error {
	out := driver.NewHtmlDriver()
	var e = echo.New()
	e.Debug = core.Context.Debug
	// Middlewares
	e.Pre(middleware.NonWWWRedirect())
	e.Pre(middleware.RemoveTrailingSlash())
	// e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	// 	return func(c echo.Context) error {
	// 		var request internal.ApiRequest
	// 		if err := c.Bind(&request) ; err != nil {
	// 			return c.NoContent(http.StatusBadRequest)
	// 		}
	// 		return next(&HttpContext { c, &ApiContext { db, request, internal.ApiResponse{} }})
	// 	}
	// })
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig {
		Skipper: middleware.DefaultSkipper,
		Format: "${time_rfc3339} ${method} ${uri} [${status}] ${error}\n",
	}))
	// Rendering
	e.Renderer = out
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusNotAcceptable
		if he, ok := err.(*echo.HTTPError) ; ok {
			code = he.Code
		}
		c.Render(code, driver.T_ERROR, struct{}{})
	}
	// Routes
	e.File("/favicon-light.ico", "html/favicon-light.ico")
	e.File("/favicon-dark.ico", "html/favicon-dark.ico")
	e.GET("/", Get)
	e.GET("/_new", GetNew)
	e.POST("/_new", PostNew)
	e.GET("/:slug", GetSlug)
	e.GET("/:slug/_edit", GetSlugEdit)
	e.POST("/:slug/_edit", PostSlugEdit)
	e.GET("/:slug/_new", GetSlugNew)
	e.POST("/:slug/_new", PostSlugNew)
	e.GET("/:slug/:progress", GetSlugProgress)
	e.GET("/:slug/_delete", GetSlugDelete)
	e.HideBanner = true
	e.HidePort = true
	addr := fmt.Sprintf("%s:%d", cli.String("ip"), cli.Int("port"))
	core.Context.Logger.Info("Listen", addr)
	e.Logger.Fatal(e.Start(addr))
	return nil
}

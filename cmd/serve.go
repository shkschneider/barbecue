package cmd

import (
	"fmt"
	"os"
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/urfave/cli/v2"
	"barbecue/driver"
	"barbecue/server"
)

var Serve = &cli.Command {
	Name: "serve",
	Usage: "serves a website",
	Action: serve,
	Flags: []cli.Flag {
		&cli.StringFlag {
			Name:  "host",
			Value: "localhost",
			Usage: "host",
		},
		&cli.IntFlag {
			Name:  "port",
			Value: 8080,
			Usage: "port",
		},
	},
}

func serve(cli *cli.Context) error {
	addr := fmt.Sprintf("%s:%d", cli.String("host"), cli.Int("port"))
	var e = echo.New()
	e.Debug = cli.String("log") == "debug" || os.Getenv("DEBUG") == "true"
	// Middlewares
	e.Pre(middleware.NonWWWRedirect())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig {
		Skipper: middleware.DefaultSkipper,
		Format: "${time_rfc3339} ${method} ${uri} [${status}] ${error}\n",
	}))
	// Rendering
	e.Renderer = driver.NewHtmlDriverRenderer()
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
	e.GET("/", server.Get)
	e.GET("/_new", server.GetNew)
	e.POST("/_new", server.PostNew)
	e.GET("/:slug", server.GetSlug)
	e.GET("/:slug/_edit", server.GetSlugEdit)
	e.POST("/:slug/_edit", server.PostSlugEdit)
	e.GET("/:slug/_new", server.GetSlugNew)
	e.POST("/:slug/_new", server.PostSlugNew)
	e.GET("/:slug/:progress", server.GetSlugProgress)
	e.GET("/:slug/_delete", server.GetSlugDelete)
	// Run
	e.HideBanner = true
	e.HidePort = true
	fmt.Println(addr)
	e.Logger.Fatal(e.Start(addr))
	return nil
}

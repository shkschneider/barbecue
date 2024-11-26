package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"net/http"
	"strings"
)

// Api

type (
	Api struct {
		echo.Echo
	}
	ApiRequest struct {
		Id			uint 	`param:"id"`
		Slug		string 	`param:"slug"`
		Title 		string 	`form:"title"`
		Description string 	`form:"description"`
		Progress	uint 	`param:"progress"`
	}
	ApiResponse struct {
		Parent		*Task	`json:"parent,omitempty"`
		Task		*Task	`json:"task"`
		Children	*[]Task `json:"children"`
	}
	ApiError struct {
		Code		int
		Message		string
	}
)

func NewApi(db *Database) (*Api, error) {
	var api Api = Api { *echo.New() }
	api.Debug = DEBUG
	// Middlewares
	api.Pre(middleware.NonWWWRedirect())
	api.Pre(middleware.RemoveTrailingSlash())
	api.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var apiRequest ApiRequest
			if err := c.Bind(&apiRequest) ; err != nil {
				return c.NoContent(http.StatusBadRequest)
			}
			return next(&LocalContext { c, apiRequest })
		}
	})
	api.Use(middleware.LoggerWithConfig(middleware.LoggerConfig {
		Skipper: middleware.DefaultSkipper,
		Format: "${time_rfc3339} ${method} ${uri} [${status}] ${error}\n",
	}))
	// Rendering
	api.Renderer = &Template {
	    templates: template.Must(template.ParseGlob(TEMPLATES)),
	}
	api.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusNotAcceptable
		if he, ok := err.(*echo.HTTPError) ; ok {
			code = he.Code
		}
		c.Render(code, tERROR, struct{}{})
	}
	// Routes
	api.File("/favicon-light.ico", "html/favicon-light.ico")
	api.File("/favicon-dark.ico", "html/favicon-dark.ico")
	api.NewRoutes(db)
	return &api, nil
}

func (api *Api) Run(addr string) {
	api.HideBanner = true
	api.HidePort = true
	log.Info(addr)
	api.Logger.Fatal(api.Start(addr))
}

// Template

const (
    TEMPLATES = "html/*.html"
	tINDEX = "index.html"
	tTASK = "task.html"
	tFORM = "form.html"
	tERROR = "error.html"
)

type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// LocalContext

type LocalContext struct {
	echo.Context
	apiRequest ApiRequest
	// apiResponse ApiResponse
}

func (c LocalContext) ok(code int, template string, data interface{}) error {
	return c.Render(code, template, data)
}

func (c LocalContext) ko(code int) error {
	return c.Render(code, tERROR, ApiError {
		Code: code,
		Message: http.StatusText(code),
	})
}

func (c LocalContext) redirect(uri string) error {
	if len(uri) == 0 || !strings.HasPrefix(uri, "/") { uri = "/" }
	return c.Redirect(http.StatusSeeOther, uri)
}

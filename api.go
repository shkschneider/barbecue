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
	api := echo.New()
	api.Debug = DEBUG
	{
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
	}
	{
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
	}
	{
		api.File("/favicon-light.ico", "html/favicon-light.ico")
		api.File("/favicon-dark.ico", "html/favicon-dark.ico")
		api.GET("/", func(_c echo.Context) error {
			c := _c.(*LocalContext)
			parents, _ := db.getParents()
			if parents == nil {
				return c.ok(http.StatusOK, tINDEX, []ApiResponse{})
			}
			var responses []ApiResponse = *new([]ApiResponse)
			for _, parent := range *parents {
				children, _ := db.getChildren(&parent)
				responses = append(responses, ApiResponse {
					Parent: nil,
					Task: &parent,
					Children: children,
				})
			}
			return c.ok(http.StatusOK, tINDEX, &responses)
		})
		api.GET("/_new", func(_c echo.Context) error {
			c := _c.(*LocalContext)
			return c.ok(http.StatusOK, tFORM, nil)
		})
		api.POST("/_new", func(_c echo.Context) error {
			c := _c.(*LocalContext)
			task, err := db.add("", c.apiRequest.Title, c.apiRequest.Description)
			if err != nil { return c.ko(http.StatusInternalServerError) }
			return c.redirect("/" + task.Slug)
		})
		api.GET("/:slug", func(_c echo.Context) error {
			c := _c.(*LocalContext)
			task, err := db.get(c.apiRequest.Slug)
			if err != nil { return c.ko(http.StatusNotFound) }
			parent, _ := db.getParent(task)
			children, _ := db.getChildren(task)
			return c.ok(http.StatusFound, tTASK, ApiResponse {
				Parent: parent,
				Task: task,
				Children: children,
			})
		})
		api.GET("/:slug/_edit", func(_c echo.Context) error {
			c := _c.(*LocalContext)
			task, err := db.get(c.apiRequest.Slug)
			if err != nil { return c.ko(http.StatusNotFound) }
			return c.ok(http.StatusFound, tFORM, ApiResponse {
				Task: task,
			})
		})
		api.POST("/:slug/_edit", func(_c echo.Context) error {
			c := _c.(*LocalContext)
			task, err := db.get(c.apiRequest.Slug)
			if err != nil { return c.ko(http.StatusNotFound) }
			task.Title = c.apiRequest.Title
			task.Description = c.apiRequest.Description
			db.update(*task)
			return c.ok(http.StatusFound, tTASK, ApiResponse {
				Task: task,
			})
		})
		api.GET("/:slug/_new", func(_c echo.Context) error {
			c := _c.(*LocalContext)
			parent, err := db.get(c.apiRequest.Slug)
			if err != nil { return c.ko(http.StatusNotFound) }
			return c.ok(http.StatusOK, tFORM, ApiResponse {
				Parent: parent,
			})
		})
		api.POST("/:slug/_new", func(_c echo.Context) error {
			c := _c.(*LocalContext)
			task, err := db.add(c.apiRequest.Slug, c.apiRequest.Title, c.apiRequest.Description)
			if err != nil { return c.ko(http.StatusNotFound) }
			return c.redirect("/" + task.Slug)
		})
		api.GET("/:slug/:progress", func(_c echo.Context) error {
			c := _c.(*LocalContext)
			task, err := db.get(c.apiRequest.Slug)
			if err != nil { return c.ko(http.StatusNotFound) }
			parent, _ := db.getParent(task)
			task.Progress = c.apiRequest.Progress
			db.update(*task)
			return c.redirect("/" + parent.Slug)
		})
		api.GET("/:slug/_delete", func(_c echo.Context) error {
			c := _c.(*LocalContext)
			task, err := db.get(c.apiRequest.Slug)
			if err != nil { return c.ko(http.StatusNotFound) }
			parent, _ := db.getParent(task)
			db.remove(*task)
			if parent != nil {
				return c.redirect("/" + parent.Slug)
			} else {
				return c.redirect("/")
			}
		})
		api.GET("/_json", func(_c echo.Context) error {
			c := _c.(*LocalContext)
			if !DEBUG {
				return c.ko(http.StatusForbidden)
			}
			tasks, _ := db.getAll()
			return c.JSON(http.StatusOK, tasks)
		})
	}
	return &Api { *api }, nil
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

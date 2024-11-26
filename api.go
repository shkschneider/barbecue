package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"net/http"
)

const (
    TEMPLATES = "html/*.html"
	tINDEX = "index.html"
	tTASK = "task.html"
	tFORM = "form.html"
	tERROR = "error.html"
)

var api *echo.Echo

type (
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
		SubTasks	*[]Task `json:"subtasks"`
	}
	Template struct {
	    templates *template.Template
	}
	LocalContext struct {
		echo.Context
		apiRequest ApiRequest
		// apiResponse ApiResponse
	}
)

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func (c LocalContext) ok(code int, template string, data interface{}) error {
	return c.Render(code, template, data)
}

func (c LocalContext) ko(code int) error {
	return c.Render(code, tERROR, struct {
		Code int
		Message string
	} {
		Code: code,
		Message: http.StatusText(code),
	})
}

func (c LocalContext) redirect(uri string) error {
	// if len(uri) == 0 || !strings.HasPrefix(uri, "/") { uri = "/" }
	return c.Redirect(http.StatusSeeOther, uri)
}

func Api() (*echo.Echo, error) {
	api = echo.New()
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
	}
	{
		api.GET("/favicon.ico", func(c echo.Context) error {
			return c.NoContent(http.StatusNoContent)
		})
		api.GET("/", func(_c echo.Context) error {
			c := _c.(*LocalContext)
			tasks, _ := GetParents()
			if tasks == nil {
				return c.ok(http.StatusOK, tINDEX, []ApiResponse{})
			}
			var responses []ApiResponse = *new([]ApiResponse)
			for _, task := range *tasks {
				subtasks, _ := GetSubTasks(task.Slug)
				responses = append(responses, ApiResponse {
					Parent: nil,
					Task: &task,
					SubTasks: subtasks,
				})
			}
			return c.ok(http.StatusOK, tINDEX, &responses)
		})
		api.GET("/new", func(_c echo.Context) error {
			c := _c.(*LocalContext)
			return c.ok(http.StatusOK, tFORM, nil)
		})
		api.POST("/new", func(_c echo.Context) error {
			c := _c.(*LocalContext)
			task, err := New("", c.apiRequest.Title, c.apiRequest.Description)
			if err != nil { return c.ko(http.StatusInternalServerError) }
			return c.redirect("/" + task.Slug)
		})
		api.GET("/:slug", func(_c echo.Context) error {
			c := _c.(*LocalContext)
			parent, _ := GetParent(c.apiRequest.Slug)
			task, err := GetTask(c.apiRequest.Slug)
			if err != nil { return c.ko(http.StatusNotFound) }
			tasks, _ := GetSubTasks(c.apiRequest.Slug)
			return c.ok(http.StatusFound, tTASK, ApiResponse {
				Parent: parent,
				Task: task,
				SubTasks: tasks,
			})
		})
		api.GET("/:slug/edit", func(_c echo.Context) error {
			c := _c.(*LocalContext)
			task, err := GetTask(c.apiRequest.Slug)
			if err != nil { return c.ko(http.StatusNotFound) }
			return c.ok(http.StatusFound, tFORM, ApiResponse {
				Task: task,
			})
		})
		api.POST("/:slug/edit", func(_c echo.Context) error {
			c := _c.(*LocalContext)
			task, err := GetTask(c.apiRequest.Slug)
			if err != nil { return c.ko(http.StatusNotFound) }
			task.Title = c.apiRequest.Title
			task.Description = c.apiRequest.Description
			Update(*task)
			return c.ok(http.StatusFound, tTASK, ApiResponse {
				Task: task,
			})
		})
		api.GET("/:slug/new", func(_c echo.Context) error {
			c := _c.(*LocalContext)
			task, err := GetTask(c.apiRequest.Slug)
			if err != nil { return c.ko(http.StatusNotFound) }
			return c.ok(http.StatusOK, tFORM, ApiResponse {
				Parent: task,
			})
		})
		api.POST("/:slug/new", func(_c echo.Context) error {
			c := _c.(*LocalContext)
			task, err := New(c.apiRequest.Slug, c.apiRequest.Title, c.apiRequest.Description)
			if err != nil { return c.ko(http.StatusNotFound) }
			return c.redirect("/" + task.Slug)
		})
		api.GET("/:slug/:progress", func(_c echo.Context) error {
			c := _c.(*LocalContext)
			parent, err := GetParent(c.apiRequest.Slug)
			task, err := GetTask(c.apiRequest.Slug)
			if err != nil { return c.ko(http.StatusNotFound) }
			task.Progress = c.apiRequest.Progress
			Update(*task)
			return c.redirect("/" + parent.Slug)
		})
		api.GET("/:slug/delete", func(_c echo.Context) error {
			c := _c.(*LocalContext)
			task, err := GetTask(c.apiRequest.Slug)
			if err != nil { return c.ko(http.StatusNotFound) }
			parent, _ := GetParent(c.apiRequest.Slug)
			Delete(*task)
			if parent != nil {
				return c.redirect("/" + parent.Slug)
			} else {
				return c.redirect("/")
			}
		})
	}
	return api, nil
}

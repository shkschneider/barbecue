package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"net/http"
	"strings"
)

var api *echo.Echo

type Request struct {
	Id			uint 	`param:"id"`
	Slug		string 	`param:"slug"`
	Title 		string 	`form:"title"`
	Description string 	`form:"description"`
	Progress	uint 	`param:"progress"`
}

type Response struct {
	Parent		*Task	`json:"parent,omitempty"`
	Task		*Task	`json:"task"`
	SubTasks	*[]Task `json:"subtasks"`
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.Map(func(r rune) rune {
		if (r < 'a' || r > 'z') && r != '-' && r != '_' && (r < '0' || r > '9') {
			return rune('-')
		} else {
			return r
		}
	}, s)
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, "--", "-")
	for strings.HasPrefix(s, "-") { s = strings.TrimPrefix(s, "-") }
	for strings.HasSuffix(s, "-") { s = strings.TrimSuffix(s, "-") }
	return s
}

type Template struct {
    templates *template.Template
}
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type apiContext struct {
	echo.Context
	Req Request
}

func Api() (*echo.Echo, error) {
	api = echo.New()
	// api.HTTPErrorHandler = func(err error, c echo.Context) {}
	api.Pre(middleware.NonWWWRedirect())
	api.Pre(middleware.RemoveTrailingSlash())
	api.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var req Request
			if err := c.Bind(&req) ; err != nil {
				return c.NoContent(http.StatusBadRequest)
			}
			return next(&apiContext { c, req })
		}
	})
	api.Use(middleware.LoggerWithConfig(middleware.LoggerConfig {
		Skipper: middleware.DefaultSkipper,
		Format: "${time_rfc3339} ${method} ${uri} [${status}] ${error}\n",
	}))
	api.Renderer = &Template {
	    templates: template.Must(template.ParseGlob("html/*.html")),
	}
	api.GET("/favicon.ico", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
	api.GET("/", func(c echo.Context) error {
		tasks, _ := GetParents()
		var responses []Response = make([]Response, 0)
		for _, task := range *tasks {
			subtasks, _ := GetSubTasks(task.Slug)
			responses = append(responses, Response {
				Parent: nil,
				Task: &task,
				SubTasks: subtasks,
			})
		}
		return c.Render(http.StatusOK, "index.html", &responses)
	})
	api.GET("/+", func(c echo.Context) error {
		return c.Render(http.StatusOK, "form.html", nil)
	})
	api.POST("/+", func(c echo.Context) error {
		req := c.(*apiContext).Req
		task, _ := New("", req.Title, req.Description)
		return c.Redirect(http.StatusSeeOther, "/" + task.Slug)
	})
	api.GET("/:slug", func(c echo.Context) error {
		req := c.(*apiContext).Req
		parent, _ := GetParent(req.Slug)
		task, _ := GetTask(req.Slug)
		tasks, _ := GetSubTasks(req.Slug)
		return c.Render(http.StatusFound, "task.html", Response {
			Parent: parent,
			Task: task,
			SubTasks: tasks,
		})
	})
	api.GET("/:slug/~", func(c echo.Context) error {
		req := c.(*apiContext).Req
		task, _ := GetTask(req.Slug)
		return c.Render(http.StatusFound, "form.html", Response {
			Task: task,
		})
	})
	api.POST("/:slug/~", func(c echo.Context) error {
		req := c.(*apiContext).Req
		task, _ := GetTask(req.Slug)
		task.Title = req.Title
		task.Description = req.Description
		Update(*task)
		return c.Render(http.StatusFound, "task.html", Response {
			Task: task,
		})
	})
	api.GET("/:slug/+", func(c echo.Context) error {
		req := c.(*apiContext).Req
		task, _ := GetTask(req.Slug)
		return c.Render(http.StatusOK, "form.html", Response {
			Task: task,
		})
	})
	api.POST("/:slug/+", func(c echo.Context) error {
		req := c.(*apiContext).Req
		task, _ := New(req.Slug, req.Title, req.Description)
		return c.Redirect(http.StatusSeeOther, "/" + task.Slug)
	})
	api.GET("/:slug/:progress", func(c echo.Context) error {
		req := c.(*apiContext).Req
		parent, _ := GetParent(req.Slug)
		task, _ := GetTask(req.Slug)
		task.Progress = req.Progress
		Update(*task)
		return c.Redirect(http.StatusSeeOther, "/" + parent.Slug)
	})
	api.GET("/:slug/-", func(c echo.Context) error {
		req := c.(*apiContext).Req
		task, _ := GetTask(req.Slug)
		Delete(*task)
		if parent, _ := GetParent(req.Slug) ; parent != nil {
			return c.Redirect(http.StatusSeeOther, "/" + parent.Slug)
		} else {
			return c.Redirect(http.StatusSeeOther, "/")
		}
	})
	return api, nil
}

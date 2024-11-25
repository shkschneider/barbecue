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

func request(c echo.Context) (*Request) {
	var request Request
	if err := c.Bind(&request) ; err != nil {
		return nil
	}
	return &request
}

func Api() (*echo.Echo, error) {
	api = echo.New()
	api.Pre(middleware.NonWWWRedirect())
	api.Pre(middleware.RemoveTrailingSlash())
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
		request := request(c) // return c.String(http.StatusBadRequest, "!")
		task, _ := New("", request.Title, request.Description)
		return c.Redirect(http.StatusSeeOther, "/" + task.Slug)
	})
	api.GET("/:slug", func(c echo.Context) error {
		request := request(c)
		parent, _ := GetParent(request.Slug)
		task, _ := GetTask(request.Slug)
		tasks, _ := GetSubTasks(request.Slug)
		return c.Render(http.StatusFound, "task.html", Response {
			Parent: parent,
			Task: task,
			SubTasks: tasks,
		})
	})
	api.GET("/:slug/~", func(c echo.Context) error {
		request := request(c)
		task, _ := GetTask(request.Slug)
		return c.Render(http.StatusFound, "form.html", Response {
			Task: task,
		})
	})
	api.POST("/:slug/~", func(c echo.Context) error {
		request := request(c)
		task, _ := GetTask(request.Slug)
		task.Title = request.Title
		task.Description = request.Description
		Update(*task)
		return c.Render(http.StatusFound, "task.html", Response {
			Task: task,
		})
	})
	api.GET("/:slug/+", func(c echo.Context) error {
		request := request(c)
		task, _ := GetTask(request.Slug)
		return c.Render(http.StatusOK, "form.html", Response {
			Task: task,
		})
	})
	api.POST("/:slug/+", func(c echo.Context) error {
		request := request(c)
		task, _ := New(request.Slug, request.Title, request.Description)
		return c.Redirect(http.StatusSeeOther, "/" + task.Slug)
	})
	api.GET("/:slug/:progress", func(c echo.Context) error {
		request := request(c)
		parent, _ := GetParent(request.Slug)
		task, _ := GetTask(request.Slug)
		task.Progress = request.Progress
		Update(*task)
		return c.Redirect(http.StatusSeeOther, "/" + parent.Slug)
	})
	api.GET("/:slug/-", func(c echo.Context) error {
		request := request(c)
		task, _ := GetTask(request.Slug)
		Delete(*task)
		if parent, _ := GetParent(request.Slug) ; parent != nil {
			return c.Redirect(http.StatusSeeOther, "/" + parent.Slug)
		} else {
			return c.Redirect(http.StatusSeeOther, "/")
		}
	})
	return api, nil
}

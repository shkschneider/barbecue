package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"net/http"
	"strings"
)

var API *echo.Echo

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

func Api() (*echo.Echo, error) {
	API = echo.New()
	API.Pre(middleware.NonWWWRedirect())
	API.Pre(middleware.RemoveTrailingSlash())
	API.Use(middleware.LoggerWithConfig(middleware.LoggerConfig {
		Skipper: middleware.DefaultSkipper,
		Format: "${time_rfc3339} ${method} ${uri} [${status}] ${error}\n",
	}))
	API.Renderer = &Template {
	    templates: template.Must(template.ParseGlob("html/*.html")),
	}
	API.GET("/favicon.ico", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
	API.GET("/", func(c echo.Context) error {
		var tasks []Task
		DB.Model(&Task{}).Where("parent_id IS NULL").Find(&tasks)
		var responses []Response = make([]Response, 0)
		for _, task := range tasks {
			if response, err := DatabaseGet(task.Slug) ; err == nil {
				responses = append(responses, *response)
			}
		}
		return c.Render(http.StatusOK, "index.html", &responses)
	})
	API.GET("/+", func(c echo.Context) error {
		return c.Render(http.StatusOK, "form.html", nil)
	})
	API.POST("/+", func(c echo.Context) error {
		var request Request ; if err := c.Bind(&request) ; err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		task, err := DatabaseSet("", request.Title, request.Description) ; if err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		return c.Redirect(http.StatusSeeOther, "/" + task.Slug)
	})
	API.GET("/:slug", func(c echo.Context) error {
		var request Request ; if err := c.Bind(&request) ; err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		response, err := DatabaseGet(request.Slug) ; if err != nil {
			return c.String(http.StatusNotFound, "!")
		}
		return c.Render(http.StatusFound, "task.html", response)
	})
	API.GET("/:slug/~", func(c echo.Context) error {
		var request Request ; if err := c.Bind(&request) ; err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		response, err := DatabaseGet(request.Slug) ; if err != nil {
			return c.String(http.StatusNotFound, "!")
		}
		return c.Render(http.StatusFound, "form.html", response)
	})
	API.POST("/:slug/~", func(c echo.Context) error {
		var request Request ; if err := c.Bind(&request) ; err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		response, err := DatabaseGet(request.Slug) ; if err != nil {
			return c.String(http.StatusNotFound, "!")
		}
		response.Task.Title = request.Title
		response.Task.Description = request.Description
		DB.Save(&response.Task)
		return c.Render(http.StatusFound, "task.html", response)
	})
	API.GET("/:slug/+", func(c echo.Context) error {
		var request Request ; if err := c.Bind(&request) ; err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		response, err := DatabaseGet(request.Slug) ; if err != nil {
			return c.String(http.StatusNotFound, "!")
		}
		return c.Render(http.StatusOK, "form.html", response)
	})
	API.POST("/:slug/+", func(c echo.Context) error {
		var request Request ; if err := c.Bind(&request) ; err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		task, err := DatabaseSet(request.Slug, request.Title, request.Description) ; if err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		return c.String(http.StatusOK, task.Slug)
	})
	API.GET("/:slug/:progress", func(c echo.Context) error {
		var request Request ; if err := c.Bind(&request) ; err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		response, err := DatabaseGet(request.Slug) ; if err != nil {
			return c.String(http.StatusNotFound, "!")
		}
		response.Task.Progress = request.Progress
		DB.Save(response.Task)
		return c.Redirect(http.StatusSeeOther, "/" + response.Parent.Slug)
	})
	API.GET("/:slug/-", func(c echo.Context) error {
		var request Request ; if err := c.Bind(&request) ; err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		response, err := DatabaseGet(request.Slug) ; if err != nil {
			return c.String(http.StatusNotFound, "!")
		}
		DB.Delete(&response.Task)
		if response.Parent != nil {
			return c.Redirect(http.StatusSeeOther, "/" + response.Parent.Slug)
		} else {
			return c.Redirect(http.StatusSeeOther, "/")
		}
	})
	return API, nil
}

package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"net/http"
	"strings"
)

type (
	EchoContext struct {
		*echo.Echo
	}
	HttpContext struct {
		echo.Context
		api 	*ApiContext
	}
)

func NewHttp(db *Database) (*EchoContext, error) {
	var e = EchoContext { echo.New() }
	e.Debug = DEBUG
	// Middlewares
	e.Pre(middleware.NonWWWRedirect())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var request ApiRequest
			if err := c.Bind(&request) ; err != nil {
				return c.NoContent(http.StatusBadRequest)
			}
			return next(&HttpContext { c, &ApiContext { db, request, ApiResponse{} }})
		}
	})
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig {
		Skipper: middleware.DefaultSkipper,
		Format: "${time_rfc3339} ${method} ${uri} [${status}] ${error}\n",
	}))
	// Rendering
	e.Renderer = &Template {
	    templates: template.Must(template.ParseGlob(TEMPLATES)),
	}
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusNotAcceptable
		if he, ok := err.(*echo.HTTPError) ; ok {
			code = he.Code
		}
		c.Render(code, tERROR, struct{}{})
	}
	// Routes
	e.File("/favicon-light.ico", "html/favicon-light.ico")
	e.File("/favicon-dark.ico", "html/favicon-dark.ico")
	//e.NewRoutes(db)
	return &e, nil
}

func (ctx *EchoContext) Run(addr string) {
	ctx.HideBanner = true
	ctx.HidePort = true
	log.Info("Listen", addr)
	ctx.Logger.Fatal(ctx.Start(addr))
}

func (ctx *HttpContext) Ok(code int, template string, data interface{}) error {
	return ctx.Render(code, template, data)
}

func (ctx *HttpContext) Ko(code int) error {
	return ctx.Render(code, tERROR, ApiError {
		Code: code,
		Message: http.StatusText(code),
	})
}

func (ctx *HttpContext) redirect(uri string) error {
	if len(uri) == 0 || !strings.HasPrefix(uri, "/") { uri = "/" }
	return ctx.Redirect(http.StatusSeeOther, uri)
}

// Main

func (ctx *EchoContext) NewRoutes(db *Database) {
	ctx.File("/favicon-light.ico", "html/favicon-light.ico")
	ctx.File("/favicon-dark.ico", "html/favicon-dark.ico")
	ctx.GET("/", func (_c echo.Context) error {
		c := _c.(*HttpContext)
		parents, _ := c.api.DB.getParents()
		if parents == nil {
			return c.Ok(http.StatusOK, tINDEX, []ApiResponse{})
		}
		var responses []ApiResponse = *new([]ApiResponse)
		for _, parent := range *parents {
			children, _ := c.api.DB.getChildren(parent)
			responses = append(responses, ApiResponse {
				Parent: nil,
				Task: &parent,
				Children: children,
			})
		}
		return c.Ok(http.StatusOK, tINDEX, &responses)
	})
	ctx.GET("/_new", func (_c echo.Context) error {
		c := _c.(*HttpContext)
		return c.Ok(http.StatusOK, tFORM, nil)
	})
	ctx.POST("/_new", func (_c echo.Context) error {
		c := _c.(*HttpContext)
		task, err := c.api.DB.add("", c.api.Request.Title, c.api.Request.Description)
		if err != nil { return c.Ko(http.StatusInternalServerError) }
		return c.redirect("/" + task.Slug)
	})
	ctx.GET("/:slug", func (_c echo.Context) error {
		c := _c.(*HttpContext)
		task, err := c.api.DB.get(c.api.Request.Slug)
		if err != nil { return c.Ko(http.StatusNotFound) }
		parent, _ := c.api.DB.getParent(*task)
		children, _ := c.api.DB.getChildren(*task)
		return c.Ok(http.StatusFound, tTASK, ApiResponse {
			Parent: parent,
			Task: task,
			Children: children,
		})
	})
	ctx.GET("/:slug/_edit", func (_c echo.Context) error {
		c := _c.(*HttpContext)
		task, err := c.api.DB.get(c.api.Request.Slug)
		if err != nil { return c.Ko(http.StatusNotFound) }
		return c.Ok(http.StatusFound, tFORM, ApiResponse {
			Task: task,
		})
	})
	ctx.POST("/:slug/_edit", func (_c echo.Context) error {
		c := _c.(*HttpContext)
		task, err := c.api.DB.get(c.api.Request.Slug)
		if err != nil { return c.Ko(http.StatusNotFound) }
		task.Title = c.api.Request.Title
		task.Description = c.api.Request.Description
		c.api.DB.update(*task)
		return c.redirect("/" + task.Slug)
	})
	ctx.GET("/:slug/_new", func (_c echo.Context) error {
		c := _c.(*HttpContext)
		parent, err := c.api.DB.get(c.api.Request.Slug)
		if err != nil { return c.Ko(http.StatusNotFound) }
		return c.Ok(http.StatusOK, tFORM, ApiResponse {
			Parent: parent,
		})
	})
	ctx.POST("/:slug/_new", func (_c echo.Context) error {
		c := _c.(*HttpContext)
		task, err := c.api.DB.add(c.api.Request.Slug, c.api.Request.Title, c.api.Request.Description)
		if err != nil { return c.Ko(http.StatusNotFound) }
		return c.redirect("/" + task.Slug)
	})
	ctx.GET("/:slug/:progress", func (_c echo.Context) error {
		c := _c.(*HttpContext)
		task, err := c.api.DB.get(c.api.Request.Slug)
		if err != nil { return c.Ko(http.StatusNotFound) }
		parent, _ := c.api.DB.getParent(*task)
		task.Progress = c.api.Request.Progress
		c.api.DB.update(*task)
		return c.redirect("/" + parent.Slug)
	})
	ctx.GET("/:slug/_delete", func (_c echo.Context) error {
		c := _c.(*HttpContext)
		task, err := c.api.DB.get(c.api.Request.Slug)
		if err != nil { return c.Ko(http.StatusNotFound) }
		parent, _ := c.api.DB.getParent(*task)
		c.api.DB.remove(*task)
		if parent != nil {
			return c.redirect("/" + parent.Slug)
		} else {
			return c.redirect("/")
		}
	})
	ctx.GET("/_json", func (_c echo.Context) error {
		c := _c.(*HttpContext)
		if !DEBUG {
			return c.Ko(http.StatusForbidden)
		}
		tasks, _ := c.api.DB.getAll()
		return c.JSON(http.StatusOK, tasks)
	})
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

func (t Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

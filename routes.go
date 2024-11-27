package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (api *Api) NewRoutes(db *Database) {
	api.File("/favicon-light.ico", "html/favicon-light.ico")
	api.File("/favicon-dark.ico", "html/favicon-dark.ico")
	// getIndex
	api.GET("/", func(_c echo.Context) error {
		c := _c.(*LocalContext)
		parents, _ := db.getParents()
		if parents == nil {
			return c.ok(http.StatusOK, tINDEX, []ApiResponse{})
		}
		var responses []ApiResponse = *new([]ApiResponse)
		for _, parent := range *parents {
			children, _ := db.getChildren(parent)
			responses = append(responses, ApiResponse {
				Parent: nil,
				Task: &parent,
				Children: children,
			})
		}
		return c.ok(http.StatusOK, tINDEX, &responses)
	})
	// getNew
	api.GET("/_new", func(_c echo.Context) error {
		c := _c.(*LocalContext)
		return c.ok(http.StatusOK, tFORM, nil)
	})
	// postNew
	api.POST("/_new", func(_c echo.Context) error {
		c := _c.(*LocalContext)
		task, err := db.add("", c.apiRequest.Title, c.apiRequest.Description)
		if err != nil { return c.ko(http.StatusInternalServerError) }
		return c.redirect("/" + task.Slug)
	})
	// getSlug
	api.GET("/:slug", func(_c echo.Context) error {
		c := _c.(*LocalContext)
		task, err := db.get(c.apiRequest.Slug)
		if err != nil { return c.ko(http.StatusNotFound) }
		parent, _ := db.getParent(*task)
		children, _ := db.getChildren(*task)
		return c.ok(http.StatusFound, tTASK, ApiResponse {
			Parent: parent,
			Task: task,
			Children: children,
		})
	})
	// getSlugEdit
	api.GET("/:slug/_edit", func(_c echo.Context) error {
		c := _c.(*LocalContext)
		task, err := db.get(c.apiRequest.Slug)
		if err != nil { return c.ko(http.StatusNotFound) }
		return c.ok(http.StatusFound, tFORM, ApiResponse {
			Task: task,
		})
	})
	// postSlugEdit
	api.POST("/:slug/_edit", func(_c echo.Context) error {
		c := _c.(*LocalContext)
		task, err := db.get(c.apiRequest.Slug)
		if err != nil { return c.ko(http.StatusNotFound) }
		task.Title = c.apiRequest.Title
		task.Description = c.apiRequest.Description
		db.update(*task)
		return c.redirect("/" + task.Slug)
	})
	// getSlugNew
	api.GET("/:slug/_new", func(_c echo.Context) error {
		c := _c.(*LocalContext)
		parent, err := db.get(c.apiRequest.Slug)
		if err != nil { return c.ko(http.StatusNotFound) }
		return c.ok(http.StatusOK, tFORM, ApiResponse {
			Parent: parent,
		})
	})
	// postSlugNew
	api.POST("/:slug/_new", func(_c echo.Context) error {
		c := _c.(*LocalContext)
		task, err := db.add(c.apiRequest.Slug, c.apiRequest.Title, c.apiRequest.Description)
		if err != nil { return c.ko(http.StatusNotFound) }
		return c.redirect("/" + task.Slug)
	})
	// getSlugProgress
	api.GET("/:slug/:progress", func(_c echo.Context) error {
		c := _c.(*LocalContext)
		task, err := db.get(c.apiRequest.Slug)
		if err != nil { return c.ko(http.StatusNotFound) }
		parent, _ := db.getParent(*task)
		task.Progress = c.apiRequest.Progress
		db.update(*task)
		return c.redirect("/" + parent.Slug)
	})
	// getSlugDelete
	api.GET("/:slug/_delete", func(_c echo.Context) error {
		c := _c.(*LocalContext)
		task, err := db.get(c.apiRequest.Slug)
		if err != nil { return c.ko(http.StatusNotFound) }
		parent, _ := db.getParent(*task)
		db.remove(*task)
		if parent != nil {
			return c.redirect("/" + parent.Slug)
		} else {
			return c.redirect("/")
		}
	})
	// getJson
	api.GET("/_json", func(_c echo.Context) error {
		c := _c.(*LocalContext)
		if !DEBUG {
			return c.ko(http.StatusForbidden)
		}
		tasks, _ := db.getAll()
		return c.JSON(http.StatusOK, tasks)
	})
}

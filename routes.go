package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// GET

func getIndex(_c echo.Context) error {
	c := _c.(*LocalContext)
	parents, _ := c.db.getParents()
	if parents == nil {
		return c.ok(http.StatusOK, tINDEX, []ApiResponse{})
	}
	var responses []ApiResponse = *new([]ApiResponse)
	for _, parent := range *parents {
		children, _ := c.db.getChildren(parent)
		responses = append(responses, ApiResponse {
			Parent: nil,
			Task: &parent,
			Children: children,
		})
	}
	return c.ok(http.StatusOK, tINDEX, &responses)
}

func getNew(_c echo.Context) error {
	c := _c.(*LocalContext)
	return c.ok(http.StatusOK, tFORM, nil)
}

func getSlug(_c echo.Context) error {
	c := _c.(*LocalContext)
	task, err := c.db.get(c.apiRequest.Slug)
	if err != nil { return c.ko(http.StatusNotFound) }
	parent, _ := c.db.getParent(*task)
	children, _ := c.db.getChildren(*task)
	return c.ok(http.StatusFound, tTASK, ApiResponse {
		Parent: parent,
		Task: task,
		Children: children,
	})
}

func getSlugEdit(_c echo.Context) error {
	c := _c.(*LocalContext)
	task, err := c.db.get(c.apiRequest.Slug)
	if err != nil { return c.ko(http.StatusNotFound) }
	return c.ok(http.StatusFound, tFORM, ApiResponse {
		Task: task,
	})
}

func getSlugNew(_c echo.Context) error {
	c := _c.(*LocalContext)
	parent, err := c.db.get(c.apiRequest.Slug)
	if err != nil { return c.ko(http.StatusNotFound) }
	return c.ok(http.StatusOK, tFORM, ApiResponse {
		Parent: parent,
	})
}

func getSlugProgress(_c echo.Context) error {
	c := _c.(*LocalContext)
	task, err := c.db.get(c.apiRequest.Slug)
	if err != nil { return c.ko(http.StatusNotFound) }
	parent, _ := c.db.getParent(*task)
	task.Progress = c.apiRequest.Progress
	c.db.update(*task)
	return c.redirect("/" + parent.Slug)
}

func getSlugDelete(_c echo.Context) error {
	c := _c.(*LocalContext)
	task, err := c.db.get(c.apiRequest.Slug)
	if err != nil { return c.ko(http.StatusNotFound) }
	parent, _ := c.db.getParent(*task)
	c.db.remove(*task)
	if parent != nil {
		return c.redirect("/" + parent.Slug)
	} else {
		return c.redirect("/")
	}
}

func getJson(_c echo.Context) error {
	c := _c.(*LocalContext)
	if !DEBUG {
		return c.ko(http.StatusForbidden)
	}
	tasks, _ := c.db.getAll()
	return c.JSON(http.StatusOK, tasks)
}

// POST

func postNew(_c echo.Context) error {
	c := _c.(*LocalContext)
	task, err := c.db.add("", c.apiRequest.Title, c.apiRequest.Description)
	if err != nil { return c.ko(http.StatusInternalServerError) }
	return c.redirect("/" + task.Slug)
}

func postSlugEdit(_c echo.Context) error {
	c := _c.(*LocalContext)
	task, err := c.db.get(c.apiRequest.Slug)
	if err != nil { return c.ko(http.StatusNotFound) }
	task.Title = c.apiRequest.Title
	task.Description = c.apiRequest.Description
	c.db.update(*task)
	return c.redirect("/" + task.Slug)
}

func postSlugNew(_c echo.Context) error {
	c := _c.(*LocalContext)
	task, err := c.db.add(c.apiRequest.Slug, c.apiRequest.Title, c.apiRequest.Description)
	if err != nil { return c.ko(http.StatusNotFound) }
	return c.redirect("/" + task.Slug)
}

// Main

func (api *Api) NewRoutes(db *Database) {
	api.File("/favicon-light.ico", "html/favicon-light.ico")
	api.File("/favicon-dark.ico", "html/favicon-dark.ico")
	api.GET("/", getIndex)
	api.GET("/_new", getNew)
	api.POST("/_new", postNew)
	api.GET("/:slug", getSlug)
	api.GET("/:slug/_edit", getSlugEdit)
	api.POST("/:slug/_edit", postSlugEdit)
	api.GET("/:slug/_new", getSlugNew)
	api.POST("/:slug/_new", postSlugNew)
	api.GET("/:slug/:progress", getSlugProgress)
	api.GET("/:slug/_delete", getSlugDelete)
	api.GET("/_json", getJson)
}

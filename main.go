package main

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
)

const NAME = "barbecue"
var VERSIONS = []string {
	"2.x",	// TODO usecase / repository
	"2.3",	// edition
	"2.2",	// progress
	"2.1",	// subtasks
	"2.0",	// golang with echo, html/template, sqlite
 	"1.0",	// html5 plain javascript, subtasks and progress
}

type Task struct {
	gorm.Model
	Slug 		string 	`param:"slug"`
	Title 		string 	`form:"title"`
	Description string 	`form:"description"`
	Progress	uint
	ParentID 	*uint
}

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

func get(db *gorm.DB, slug string) (*Response, error) {
	var task Task
	result := db.Model(&Task{}).Where(Task { Slug: slug }).First(&task)
	var parent *Task = new(Task)
	db.Model(&Task{}).First(parent, task.ParentID)
	if parent.ID == 0 { parent = nil }
	var tasks *[]Task = new([]Task)
	db.Model(&Task{}).Where(Task { ParentID: &task.ID }).Order("progress").Find(tasks)
	if len(*tasks) == 0 {
		tasks = nil
	} else {
		var p uint = 0 ; for _, t := range *tasks { p += t.Progress }
		task.Progress = (p / uint(len(*tasks)))
		db.Save(&task)
	}
	return &Response {
		Parent: parent,
		Task: &task,
		SubTasks: tasks,
	}, result.Error
}

func set(db *gorm.DB, slug string, title string, description string) (Task, error) {
	var response *Response
	if len(slug) > 0 {
		if r, err := get(db, slug) ; err != nil {
			return Task{}, err
		} else {
			response = r
		}
	} else {
		response = &Response{}
	}
	task := Task {
		Slug: slugify(title),
		Title: title,
		Description: description,
		ParentID: nil,
	}
	if response.Task != nil {
		task.Slug = fmt.Sprintf("%v-%s", response.Task.ID, task.Slug)
		task.ParentID = &response.Task.ID
	}
	result := db.Create(&task)
	return task, result.Error
}

type Template struct {
    templates *template.Template
}
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// database
	db, err := gorm.Open(sqlite.Open(NAME + ".sqlite"), &gorm.Config{}) ; if err != nil { panic("database") }
	db.AutoMigrate(&Task{})
	if os.Getenv("DEBUG") == "true" {
		fmt.Println(slugify("This #Is_A_Slugify Test!!!"))
		fmt.Println(slug.Make("This #Is_A_Slugi Test!!!"))
		db.Session(&gorm.Session { AllowGlobalUpdate: true }).Delete(&Task{})
		set(db, "", "Title1", "Description _One_")
		set(db, "", "Title2", "Description _Two_")
		set(db, "", "Title3", "Description _Three_")
		set(db, "title1", "Title11", "Description11")
		set(db, "title1", "Title12", "Description12")
		set(db, "title1", "Title13", "Description13")
		var tasks []Task
		db.Model(&Task{}).Find(&tasks)
		for _, task := range tasks {
			fmt.Println(task.ID, task.Slug)
		}
	}
	// router
	e := echo.New()
	e.Pre(middleware.NonWWWRedirect())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig {
		Skipper: middleware.DefaultSkipper,
		Format: "${time_rfc3339} ${method} ${uri} [${status}] ${error}\n",
	}))
	// renderer
	e.Renderer = &Template {
	    templates: template.Must(template.ParseGlob("html/*.html")),
	}
	// routes
	e.GET("/favicon.ico", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
	e.GET("/", func(c echo.Context) error {
		var tasks []Task
		db.Model(&Task{}).Where("parent_id IS NULL").Find(&tasks)
		var responses []Response = make([]Response, 0)
		for _, task := range tasks {
			if response, err := get(db, task.Slug) ; err == nil {
				responses = append(responses, *response)
			}
		}
		return c.Render(http.StatusOK, "index.html", &responses)
	})
	e.GET("/+", func(c echo.Context) error {
		return c.Render(http.StatusOK, "form.html", nil)
	})
	e.POST("/+", func(c echo.Context) error {
		var request Request ; if err := c.Bind(&request) ; err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		task, err := set(db, "", request.Title, request.Description) ; if err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		return c.Redirect(http.StatusSeeOther, "/" + task.Slug)
	})
	e.GET("/:slug", func(c echo.Context) error {
		var request Request ; if err := c.Bind(&request) ; err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		response, err := get(db, request.Slug) ; if err != nil {
			return c.String(http.StatusNotFound, "!")
		}
		return c.Render(http.StatusFound, "task.html", response)
	})
	e.GET("/:slug/~", func(c echo.Context) error {
		var request Request ; if err := c.Bind(&request) ; err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		response, err := get(db, request.Slug) ; if err != nil {
			return c.String(http.StatusNotFound, "!")
		}
		return c.Render(http.StatusFound, "form.html", response)
	})
	e.POST("/:slug/~", func(c echo.Context) error {
		var request Request ; if err := c.Bind(&request) ; err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		response, err := get(db, request.Slug) ; if err != nil {
			return c.String(http.StatusNotFound, "!")
		}
		response.Task.Title = request.Title
		response.Task.Description = request.Description
		db.Save(&response.Task)
		return c.Render(http.StatusFound, "task.html", response)
	})
	e.GET("/:slug/+", func(c echo.Context) error {
		var request Request ; if err := c.Bind(&request) ; err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		response, err := get(db, request.Slug) ; if err != nil {
			return c.String(http.StatusNotFound, "!")
		}
		return c.Render(http.StatusOK, "form.html", response)
	})
	e.GET("/:slug/:progress", func(c echo.Context) error {
		var request Request ; if err := c.Bind(&request) ; err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		response, err := get(db, request.Slug) ; if err != nil {
			return c.String(http.StatusNotFound, "!")
		}
		response.Task.Progress = request.Progress
		db.Save(response.Task)
		return c.Redirect(http.StatusSeeOther, "/" + response.Parent.Slug)
	})
	e.GET("/:slug/-", func(c echo.Context) error {
		var request Request ; if err := c.Bind(&request) ; err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		response, err := get(db, request.Slug) ; if err != nil {
			return c.String(http.StatusNotFound, "!")
		}
		db.Delete(&response.Task)
		if response.Parent != nil {
			return c.Redirect(http.StatusSeeOther, "/" + response.Parent.Slug)
		} else {
			return c.Redirect(http.StatusSeeOther, "/")
		}
	})
	e.POST("/:slug/+", func(c echo.Context) error {
		var request Request ; if err := c.Bind(&request) ; err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		task, err := set(db, request.Slug, request.Title, request.Description) ; if err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		return c.String(http.StatusOK, task.Slug)
	})
	e.Logger.Fatal(e.Start("0.0.0.0:8080"))
}

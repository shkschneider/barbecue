package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"io"
	"net/http"
	"strings"
	"unicode"
)

const NAME = "barbecue"
var VERSIONS = []string {
	"2.1",	// TODO usecase, repository
	"2.0",	// golang with echo, html/template, sqlite
 	"1.0",	// html5 plain javascript, subtasks and progress
}

type Task struct {
	gorm.Model
	Slug 		string
	Title 		string 	`form:"title"`
	Description string 	`form:"description"`
	ParentID 	*uint 	`gorm:"foreignkey:ID"`
	Parent 		*Task 	`gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func slugify(s string) string {
	return strings.Map(func(r rune) rune {
		if r > unicode.MaxASCII { return -1 } else { return r }
	}, s)
}

func get(db *gorm.DB, slug string) (Task, error) {
	var task Task
	result := db.Model(&Task{}).Where("Slug = ?", slug).First(&task)
	return task, result.Error
}

func set(db *gorm.DB, slug string, title string, description string) (Task, error) {
	parent, err := get(db, slug) ; if err != nil {
		return Task{}, err
	}
	task := Task {
		Slug: slugify(title),
		Title: title,
		Description: description,
		ParentID: &parent.ID,
		Parent: &parent,
	}
	db.Create(&task)
	return task, nil
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
	db.Session(&gorm.Session { AllowGlobalUpdate: true }).Delete(&Task{})
	t1 := Task { Slug: "title1", Title: "Title1", Description: "Description1" }
	db.Create(&t1)
	db.Create(&Task { Slug: "title2", Title: "Title2", Description: "Description2" })
	db.Create(&Task { Slug: "title3", Title: "Title3", Description: "Description3" })
	db.Create(&Task { Slug: "title11", Title: "Title11", Description: "Title11", Parent: &t1 })
	db.Create(&Task { Slug: "title12", Title: "Title12", Description: "Title12", Parent: &t1 })
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
	e.GET("/", func(c echo.Context) error {
		var tasks []Task
		db.Model(&Task{}).Where("parent_id IS NULL").Find(&tasks)
		return c.Render(http.StatusOK, "index.html", tasks)
	})
	e.GET("/:slug", func(c echo.Context) error {
		task, err := get(db, c.Param("slug")) ; if err != nil {
			return c.String(http.StatusNotFound, "!")
		}
		return c.Render(http.StatusFound, "task.html", task)
	})
	e.GET("/+", func(c echo.Context) error {
		return c.Render(http.StatusOK, "form.html", nil)
	})
	e.GET("/:slug/+", func(c echo.Context) error {
		task, err := get(db, c.Param("slug"))
		if err != nil {
			return c.String(http.StatusNotFound, "!")
		}
		return c.Render(http.StatusOK, "form.html", task)
	})
	e.POST("/+", func(c echo.Context) error {
		task, err := set(db, "", c.FormValue("title"), c.FormValue("description"))
		if err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		return c.Redirect(http.StatusOK, "/" + task.Slug)
	})
	e.POST("/:slug/+", func(c echo.Context) error {
		task, err := set(db, c.Param("slug"), c.FormValue("title"), c.FormValue("description"))
		if err != nil {
			return c.String(http.StatusBadRequest, "!")
		}
		return c.String(http.StatusOK, task.Slug)
	})
	e.Logger.Fatal(e.Start("0.0.0.0:8080"))
}

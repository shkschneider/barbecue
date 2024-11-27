package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

const NAME = "barbecue"
var VERSIONS = []string {
	"2.8", // rc
	"2.7", // polish
	"2.6", // refactor
	"2.5", // no globals
	"2.4", // api, database
	"2.3", // edition
	"2.2", // progress
	"2.1", // subtasks
	"2.0", // golang with echo, html/template, sqlite
 	"1.0", // html5 plain javascript, subtasks and progress
}
var DEBUG bool = false

type Task struct {
	gorm.Model			`json:"-"`
	Slug 		string 	`json:"slug",param:"slug"`
	Title 		string 	`json:"title",form:"title"`
	Description string 	`json:"description",form:"description"`
	Progress	uint	`json:"progress",form:"progress"`
	Super 		*uint	`json:"-"`
}

func main() {
	DEBUG = (os.Getenv("DEBUG") == "true")
	if DEBUG {
		NewLog(os.Stderr, LogLevelDebug, true)
	} else {
		NewLog(os.Stderr, LogLevelInfo, true)
	}
	log.Info("Database...")
	db, err := NewDatabase(sqlite.Open(NAME + ".sqlite"))
	if err != nil {
		log.Panic("database")
	}
	log.Info("Api...")
	api, err := NewApi(db)
	if err != nil {
		log.Panic("api")
	}
	api.Run("0.0.0.0:8080")
}

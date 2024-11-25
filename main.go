package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const NAME = "barbecue"
var VERSIONS = []string {
	"2.5",	// no globals
	"2.4",	// api, database
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

func main() {
	if _, err := Database(sqlite.Open(NAME + ".sqlite")) ; err != nil {
		panic("repository")
	}
	if api, err := Api() ; err != nil {
		panic("usecase")
	} else {
		api.Logger.Fatal(api.Start("0.0.0.0:8080"))
	}
}

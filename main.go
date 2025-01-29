package main

import (
	"gorm.io/driver/sqlite"
	"os"
	"github.com/urfave/cli/v2"
	"barbecue/cmd"
	"barbecue/core"
	"barbecue/data"
)

const NAME = "barbecue"
var VERSIONS = []string {
	"4.0", // modularized
	"3.0", // cli
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

func main() {
	core.Context.Debug = (os.Getenv("DEBUG") == "true")
	if core.Context.Debug {
		core.Context.Logger = core.NewLogger(os.Stderr, core.LevelDebug)
	} else {
		core.Context.Logger = core.NewLogger(os.Stderr, core.LevelInfo)
	}
	if db, err := data.New(sqlite.Open(NAME + ".sqlite"), false) ; err != nil {
		core.Context.Logger.Panic(err)
	} else {
		core.Context.Database = db
	}
	(&cli.App {
		Name: "barbecue",
		Usage: "...",
		Commands: []*cli.Command {
			cmd.Add,
			cmd.List,
			cmd.Remove,
			cmd.Serve,
			cmd.Progress,
		},
	}).Run(os.Args)
}

package main

import (
	"fmt"
	"os"
	"time"
	"github.com/urfave/cli/v2"
	"barbecue/cmd"
	"barbecue/core"
	"barbecue/data"
)

const NAME = "barbecue"
var VERSIONS = []string {
	"4.5", // rc with go-task
	"4.4", // log levels
	"4.3", // db location
	"4.2", // errors
	"4.1", // drivers
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
	commands := []*cli.Command {
		// cli
		cmd.Add,
		cmd.List,
		cmd.Progress,
		cmd.Remove,
		cmd.Show,
		// http
		cmd.Serve,
	}
	flags := []cli.Flag{
		&cli.StringFlag {
			Name: "database",
			Value: "~/.local/state/" + NAME + ".sqlite",
		},
		&cli.StringFlag { Name: "log", Value: "error" },
	}
	before := func(cli *cli.Context) error {
		// log
		switch log := cli.String("log") ; log {
			case "debug":
				core.Log = core.NewLogger(os.Stderr, core.LevelDebug)
			case "info":
				core.Log = core.NewLogger(os.Stderr, core.LevelInfo)
			case "warning":
				core.Log = core.NewLogger(os.Stderr, core.LevelWarning)
			case "error":
				core.Log = core.NewLogger(os.Stderr, core.LevelError)
			default:
				core.Log = core.NewLogger(os.Stderr, core.LevelPanic)
		}
		core.Log.Info(fmt.Sprintf("%s v%s", NAME, VERSIONS[0]))
		// db
		var err error
		core.Database, err = data.NewDatabase(cli.String("database"), NAME)
		if err != nil {
			core.Log.Panic(err)
			return err
		}
		core.Log.Info(core.Database.Path)
		return nil
	}
	(&cli.App {
		Name: NAME,
		Version: "v" + VERSIONS[0],
        Compiled: time.Now(),
		Usage: "easy tasks and subtasks !",
		Commands: commands,
		Flags: flags,
		Before: before,
	}).Run(os.Args)
}

package cmd

import (
	"errors"
	"github.com/urfave/cli/v2"
	"barbecue/api"
	"barbecue/core"
	"barbecue/data"
	"barbecue/driver"
)

var Remove = &cli.Command {
	Name: "remove",
	Usage: "...",
	ArgsUsage: "<idOrSlug>",
	Action: remove,
	Flags: []cli.Flag {
		&cli.BoolFlag {
			Name:  "children",
			Value: false,
			Usage: "...",
		},
	},
}

func removeRecursive(task data.Task) {
	out := driver.NewStdoutDriver()
	children, err := api.GetChildren(task)
	if err == nil {
		for _, child := range *children {
			out.Output(&child)
		}
	}
}

func remove(cli *cli.Context) error {
	out := driver.NewStdoutDriver()
	tasks, err := api.GetByIdOrSlug(cli.Args().Get(0))
	if err != nil {
		core.Context.Logger.Error("not found")
		return err
	}
	if tasks == nil || len(*tasks) == 0 {
		core.Context.Logger.Error("nothing to remove")
		return errors.New("nothing to remove")
	}
	if len(*tasks) > 1 {
		for _, task := range *tasks {
			out.Output(&task)
		}
		core.Context.Logger.Warning("too ambiguous")
		return nil
	}
	task := (*tasks)[0]
	out.Output(&task)
	if cli.Bool("children") == true {
		api.RemoveRecursive(task)
	} else {
		api.Remove(task)
	}
	return nil
}

package cmd

import (
	"github.com/urfave/cli/v2"
	"barbecue/api"
	"barbecue/core"
	"barbecue/driver"
)

var Add = &cli.Command {
	Name: "add",
	Usage: "adds a task",
	Action: add,
	ArgsUsage: "<title> [description]",
	Flags: []cli.Flag {
		&cli.StringFlag {
			Name:  "parent",
			Aliases: []string { "p" },
			Value: "",
			Usage: "<idOrSlug>",
		},
	},
}

func add(cli *cli.Context) error {
	out := driver.NewStdoutDriver()
	title := cli.Args().Get(0)
	description := cli.Args().Get(1)
	tasks, err := api.Add(title, description)
	if err != nil {
		core.Log.Error("Add", err)
		return err
	}
	task := (*tasks)[0]
	if len(cli.String("parent")) > 0 {
		parents, err := api.GetByIdOrSlug(cli.String("parent"))
		if err != nil {
			core.Log.Error("Add", err)
			return err
		}
		parent := (*parents)[0]
		task.Super = &parent.ID
		core.Database.Update(task)
	}
	out.Out(driver.NewStdoutDriverData(task))
	return nil
}

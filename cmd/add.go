package cmd

import (
	"errors"
	"github.com/urfave/cli/v2"
	"barbecue/api"
	"barbecue/core"
	"barbecue/driver"
)

var Add = &cli.Command {
	Name: "add",
	Usage: "...",
	Action: add,
	ArgsUsage: "<title> [description]",
	Flags: []cli.Flag {
		&cli.StringFlag {
			Name:  "parent",
			Aliases: []string { "p" },
			Value: "",
			Usage: "...",
		},
	},
}

func add(cli *cli.Context) error {
	out := driver.NewStdoutDriver()
	tasks, err := api.Add(cli.Args().Get(0), cli.Args().Get(1))
	if err != nil {
		return err
	}
	if len(cli.String("parent")) > 0 {
		parents, err := api.GetByIdOrSlug(cli.String("parent"))
		if err == nil {
			(*tasks)[0].Super = &((*parents)[0].ID)
			core.Context.Database.Update((*tasks)[0])
		}
	}
	out.Output(&(*tasks)[0])
	return errors.New("<title> [description]")
}

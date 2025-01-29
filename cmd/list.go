package cmd

import (
	"github.com/urfave/cli/v2"
	"barbecue/api"
	"barbecue/driver"
)

var List = &cli.Command {
	Name: "list",
	Usage: "...",
	Action: list,
	ArgsUsage: "[isOrSlug ...]",
	Flags: []cli.Flag {
		&cli.BoolFlag {
			Name:  "recursive",
			Aliases: []string { "r" },
			Value: false,
			Usage: "...",
		},
	},
}

func list(cli *cli.Context) error {
	out := driver.NewStdoutDriver()
	if cli.NArg() == 0 {
		tasks, err := api.GetParents()
		for _, task := range *tasks {
			if task.Progress < 100 && cli.Bool("recursive") {
				out.OutputRecursive(0, &task)
			} else {
				out.Output(&task)
			}
		}
		return err
	} else {
		tasks, err := api.GetByIdOrSlug(cli.Args().Get(0))
		for _, task := range *tasks {
			if cli.Bool("recursive") {
				out.OutputRecursive(0, &task)
			} else {
				out.Output(&task)
			}
		}
		return err
	}
	return nil
}

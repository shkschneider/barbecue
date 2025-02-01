package cmd

import (
	"github.com/urfave/cli/v2"
	"barbecue/api"
	"barbecue/core"
	"barbecue/driver"
)

var Show = &cli.Command {
	Name: "show",
	Usage: "shows a single task (and subtasks)",
	Action: show,
	ArgsUsage: "<isOrSlug>",
	Flags: []cli.Flag{},
}

func show(cli *cli.Context) error {
	out := driver.NewStdoutDriver()
	idOrSlug := cli.Args().Get(0)
	tasks, err := api.GetByIdOrSlug(idOrSlug)
	if err != nil {
		core.Log.Error("Show", core.ErrNotFound)
		return core.ErrNotFound
	} else if tasks == nil {
		core.Log.Error("Show", core.ErrNothing)
		return core.ErrNothing
	}
	for _, task := range *tasks {
		out.Out(driver.NewStdoutDriverData(task))
	}
	return err
}

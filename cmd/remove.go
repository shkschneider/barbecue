package cmd

import (
	"github.com/urfave/cli/v2"
	"barbecue/api"
	"barbecue/core"
	"barbecue/driver"
)

var Remove = &cli.Command {
	Name: "remove",
	Usage: "removes a task (and subtasks!)",
	ArgsUsage: "<idOrSlug>",
	Action: remove,
	Flags: []cli.Flag{},
}

func remove(cli *cli.Context) error {
	out := driver.NewStdoutDriver()
	tasks, err := api.GetByIdOrSlug(cli.Args().Get(0))
	if err != nil {
		core.Log.Error("Remove", err)
		return err
	} else if tasks == nil {
		core.Log.Error("Remove", core.ErrNothing)
		return core.ErrNothing
	}
	if len(*tasks) > 1 {
		for _, task := range *tasks {
			out.Out(driver.NewStdoutDriverData(task))
		}
		core.Log.Error("Remove", core.ErrAmbiguous)
		return core.ErrAmbiguous
	}
	task := (*tasks)[0]
	api.RemoveRecursive(task)
	out.Out(driver.NewStdoutDriverData(task))
	return nil
}

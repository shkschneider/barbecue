package cmd

import (
	"strconv"
	"github.com/urfave/cli/v2"
	"barbecue/api"
	"barbecue/core"
	"barbecue/driver"
)

var Progress = &cli.Command {
	Name: "progress",
	Usage: "...",
	Action: progress,
	ArgsUsage: "<idOrSlug> <%>",
	Flags: []cli.Flag{},
}

func progress(cli *cli.Context) error {
	out := driver.NewStdoutDriver()
	tasks, err := api.GetByIdOrSlug(cli.Args().Get(0))
	if err != nil {
		return err
	}
	task := (*tasks)[0]
	pc, err := strconv.ParseUint(cli.Args().Get(1), 10, 16)
	if err != nil {
		core.Context.Logger.Error(err)
		return err
	}
	task.Progress = uint(pc)
	if _, err := api.Update(task) ; err != nil {
		core.Context.Logger.Error(err)
		return err
	}
	out.Output(&task)
	return nil
}

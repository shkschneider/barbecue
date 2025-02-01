package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"barbecue/api"
	"barbecue/core"
	"barbecue/driver"
)

var List = &cli.Command {
	Name: "list",
	Usage: "lists all tasks (recursively)",
	Action: list,
	ArgsUsage: "",
	Flags: []cli.Flag{},
}

func list(cli *cli.Context) error {
	out := driver.NewStdoutDriver()
	tasks, err := api.GetParents()
	if err != nil {
		core.Log.Error("List", err)
		return err
	} else if tasks == nil {
		core.Log.Error("List", core.ErrNothing)
		return core.ErrNothing
	}
	out.Out(driver.NewStdoutDriverData((*tasks)...))
	core.Log.Info(fmt.Sprintf("%d task(s)", len(*tasks)))
	return nil
}

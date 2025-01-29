package driver

import (
	"fmt"
	"strings"
	"barbecue/core"
	"barbecue/data"
)

type StdoutDriver struct {
	core.Driver
}

func NewStdoutDriver() *StdoutDriver {
	return &StdoutDriver{}
}

func output(depth uint, task *data.Task) {
	var p rune
	if task.Progress == 0 {
		p = 'x'
	} else if task.Progress == 100 {
		p = '✓'
	} else {
		p = '…'
	}
	var h string
	if depth == 0 {
		h = ""
	} else {
		h = fmt.Sprintf("%s- ", strings.Repeat(" ", int(depth - 1) * 2))
	}
	fmt.Printf("%s#%d [%c %d%%] '%s'\n", h, task.ID, p, task.Progress, task.Slug)
	// fmt.Printf(
	// 	"%s#%v '%s' %d%% \"%s\" \"\"\"%s\"\"\"\n",
	// 	strings.Repeat("-", int(depth)),
	// 	task.ID,
	// 	task.Slug,
	// 	task.Progress,
	// 	task.Title,
	// 	task.Description,
	// )
}

func (d *StdoutDriver) Output(task *data.Task) {
	output(0, task)
}

func (d *StdoutDriver) OutputRecursive(depth uint, task *data.Task) {
	output(depth, task)
	tasks, err := core.Context.Database.GetChildren(*task)
	if err != nil || tasks == nil { return }
	for _, task := range *tasks {
		d.OutputRecursive(depth + 1, &task)
	}
}

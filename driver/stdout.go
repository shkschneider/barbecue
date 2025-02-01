package driver

import (
	"fmt"
	"strings"
	"barbecue/core"
	"barbecue/data"
)

type StdoutDriver struct {}

func NewStdoutDriver() core.Driver {
	return StdoutDriver{}
}

func single(depth uint, task *data.Task) {
	var h string
	if depth == 0 {
		h = ""
	} else {
		h = strings.Repeat(" ", int(depth))
	}
	var p string
	if task.Progress == 0 {
		p = "x"
	} else if task.Progress == 100 {
		p = "✓"
	} else {
		p = "…"
	}
	fmt.Printf("\t%s%s %d %d%% '%s'\n", h, p, task.ID, task.Progress, task.Slug)
}

func recursive(depth uint, task *data.Task) {
	single(depth, task)
	tasks, err := core.Database.GetChildren(*task)
	if err != nil || tasks == nil { return }
	for _, task := range *tasks {
		if depth == 0 {
			recursive(uint(len(fmt.Sprintf("%v", *task.Super)) + 1), &task)
		} else {
			recursive(depth + uint(len(fmt.Sprintf("%v", *task.Super)) + 1), &task)
		}
	}
}

func (d StdoutDriver) Out(data interface{}) error {
	tasks := data.(StdoutDriverData).tasks
	if len(tasks) == 0 { return nil }
	fmt.Println("")
	for i, task := range tasks {
		if i > 0 {
			fmt.Println("")
		}
		recursive(0, &task)
	}
	fmt.Println("")
	return nil
}

func (d StdoutDriver) Err(code int, msg string) error {
	core.Log.Debug("ERR", code, msg)
	return nil
}

// Data

type StdoutDriverData struct {
	tasks	[]data.Task
}

func NewStdoutDriverData(tasks ...data.Task) StdoutDriverData {
	return StdoutDriverData { *data.NewTasks(tasks...) }
}

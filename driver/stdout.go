package driver

import (
	"fmt"
	"strings"
	"barbecue/core"
	"barbecue/data"

	"github.com/fatih/color"
)

type StdoutDriver struct {}

func NewStdoutDriver() core.Driver {
	return StdoutDriver{}
}

func single(depth uint, task *data.Task) {
	var fg *color.Color
	fmt.Printf("\t")
	fmt.Printf(strings.Repeat(" ", int(depth * 2)))
	{
		if task.Progress == 0 {
			fg = color.New(color.FgRed, color.Bold)
			fg.Print("x")
		} else if task.Progress == 100 {
			fg = color.New(color.FgGreen, color.Bold)
			fg.Print("✓")
		} else {
			fg = color.New(color.FgYellow, color.Bold)
			fg.Print("…")
		}
	}
	fmt.Print(" ")
	{
		fmt.Print("[")
		pc := int(task.Progress) / 10
		fg.Print(strings.Repeat("-", pc))
		fmt.Print(strings.Repeat(" ", 3 - len(fmt.Sprintf("%d", int(task.Progress)))))
		fmt.Print(strings.Repeat(" ", 10 - pc))
		fg.Print(int(task.Progress))
		fmt.Print("%")
		fmt.Print("]")
	}
	fmt.Print(" ")
	{
		fmt.Print("#")
		fg.Printf("%d", task.ID)
	}
	fmt.Print(" ")
	{
		fmt.Print("'")
		color.New(color.FgBlue, color.Bold).Printf("%s", task.Slug)
		fmt.Print("'")
	}
	fmt.Print(" ")
	fmt.Print("\n")
}

func recursive(depth uint, task *data.Task) {
	single(depth, task)
	tasks, err := core.Database.GetChildren(*task)
	if err != nil || tasks == nil { return }
	for _, task := range *tasks {
		// recursive(depth + uint(len(fmt.Sprintf("%v", *task.Super)) + 1), &task)
		recursive(depth + 1, &task)
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

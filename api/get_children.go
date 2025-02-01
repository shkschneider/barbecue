package api

import (
	"barbecue/core"
	"barbecue/data"
)

func GetChildren(task data.Task) (*[]data.Task, error) {
	children, err := core.Database.GetChildren(task)
	if err != nil {
		core.Log.Debug("GetChildren", err)
		return nil, err
	} else if children == nil || len(*children) == 0 {
		core.Log.Debug("GetChildren", "nothing")
		return data.NewTasks(), nil
	}
	var tasks []data.Task
	for _, child := range *children {
		task, err := GetById(child.ID)
		if err != nil {
			core.Log.Debug("GetChildren", err)
		} else if task == nil {
			core.Log.Debug("GetChildren", "nothing")
			continue
		}
		tasks = append(tasks, *task)
	}
	return &tasks, nil
}

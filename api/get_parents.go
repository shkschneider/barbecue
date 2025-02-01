package api

import (
	"barbecue/core"
	"barbecue/data"
)

func GetParents() (*[]data.Task, error) {
	parents, err := core.Database.GetParents()
	if err != nil {
		core.Log.Debug("GetParents", err)
		return nil, err
	} else if parents == nil || len(*parents) == 0 {
		core.Log.Debug("GetParents", "nothing")
		return data.NewTasks(), nil
	}
	var tasks []data.Task
	for _, parent := range *parents {
		task, err := GetById(parent.ID)
		if err != nil {
			core.Log.Debug("GetParents", err)
		} else if task == nil {
			core.Log.Debug("GetParents", "nothing")
		} else {
			tasks = append(tasks, *task)
		}
	}
	return &tasks, nil
}

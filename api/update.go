package api

import (
	"barbecue/core"
	"barbecue/data"
)

func Update(task data.Task) ([]data.Task, error) {
	err := core.Database.Update(task)
	if err != nil {
		core.Log.Debug("Update", err)
		return nil, err
	}
	return *data.NewTasks(task), nil
}

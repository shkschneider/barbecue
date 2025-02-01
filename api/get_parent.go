package api

import (
	"barbecue/core"
	"barbecue/data"
)

func GetParent(task data.Task) (*data.Task, error) {
	parent, err := core.Database.GetParent(task)
	if err != nil {
		core.Log.Debug("GetParent", err)
		return nil, err
	} else if parent == nil {
		core.Log.Debug("GetParent", "nothing")
	}
	return parent, nil
}

package api

import (
	"errors"
	"barbecue/core"
	"barbecue/data"
)

func GetById(id uint) (*data.Task, error) {
	task, err := core.Database.GetById(id)
	if err != nil {
		core.Log.Debug("GetById", err)
		return nil, err
	} else if task == nil {
		core.Log.Debug("GetById", "not found")
		return nil, errors.New("not found")
	}
	return task, nil
}

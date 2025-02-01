package api

import (
	"errors"
	"barbecue/core"
	"barbecue/data"
)

func GetBySlug(slug string) (*[]data.Task, error) {
	tasks, err := core.Database.GetBySlug(data.Slugify(slug))
	if err != nil {
		core.Log.Debug("GetBySlug", err)
		return nil, err
	} else if tasks == nil || len(*tasks) == 0 {
		core.Log.Debug("GetBySlug", "not found")
		return nil, errors.New("not found")
	}
	return tasks, nil
}

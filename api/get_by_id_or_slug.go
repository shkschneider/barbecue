package api

import (
	"errors"
	"strconv"
	"barbecue/core"
	"barbecue/data"
)

func GetByIdOrSlug(idOrSlug string) (*[]data.Task, error) {
	if id, err := strconv.ParseUint(idOrSlug, 10, 16) ; err == nil {
		task, err := core.Database.GetById(uint(id))
		if err != nil {
			core.Log.Debug("GetByIdOrSlug", err)
			return nil, err
		} else if task == nil {
			core.Log.Debug("GetByIdOrSlug", "Nothing")
			return nil, errors.New("Nothing")
		}
		return data.NewTasks(*task), nil
	} else {
		slug := data.Slugify(idOrSlug)
		tasks, err := core.Database.GetBySlug(slug)
		if err != nil {
			core.Log.Debug("GetByIdOrSlug", err)
			return nil, err
		} else if tasks == nil || len(*tasks) == 0 {
			core.Log.Debug("GetByIdOrSlug", "Nothing")
			return nil, errors.New("Nothing")
		}
		return tasks, nil
	}
}

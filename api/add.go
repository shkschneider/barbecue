package api

import (
	"barbecue/core"
	"barbecue/data"
)

func Add(title string, description string) (*[]data.Task, error) {
	slug := data.Slugify(title)
	core.Database.Insert(slug, title, description)
	tasks, err := GetBySlug(slug)
	if err != nil {
		core.Log.Debug("Add", err)
		return nil, err
	} else if tasks == nil {
		return data.NewTasks(), nil
	}
	return tasks, nil
}

package api

import (
	"strconv"
	"barbecue/core"
	"barbecue/data"
)

func GetById(id uint) (*data.Task, error) {
	task, err := core.Context.Database.GetById(id)
	return task, err
}

func GetBySlug(slug string) (*[]data.Task, error) {
	tasks, err := core.Context.Database.GetBySlug(data.Slugify(slug))
	return tasks, err
}

func GetByIdOrSlug(idOrSlug string) (*[]data.Task, error) {
	if id, err := strconv.ParseUint(idOrSlug, 10, 16) ; err == nil {
		task, err := core.Context.Database.GetById(uint(id))
		if err != nil {
			return nil, err
		}
		return &[]data.Task { *task }, nil
	} else {
		slug := data.Slugify(idOrSlug)
		tasks, err := core.Context.Database.GetBySlug(slug)
		if err != nil {
			return nil, err
		}
		return tasks, nil
	}
}

func GetParents() (*[]data.Task, error) {
	parents, err := core.Context.Database.GetParents()
	if err != nil {
		return nil, err
	}
	var tasks = *new([]data.Task)
	for _, parent := range *parents {
		task, _ := GetById(parent.ID)
		tasks = append(tasks, *task)
	}
	return &tasks, nil
}

func GetChildren(task data.Task) (*[]data.Task, error) {
	children, err := core.Context.Database.GetChildren(task)
	if err != nil {
		return nil, err
	}
	var tasks = *new([]data.Task)
	for _, child := range *children {
		task, _ := GetById(child.ID)
		tasks = append(tasks, *task)
	}
	return &tasks, nil
}

func GetParent(task data.Task) (*data.Task, error) {
	return core.Context.Database.GetParent(task)
}

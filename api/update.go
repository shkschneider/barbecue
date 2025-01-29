package api

import (
	"barbecue/core"
	"barbecue/data"
)

func Update(task data.Task) (*[]data.Task, error) {
	err := core.Context.Database.Update(task)
	return &[]data.Task { task }, err
}

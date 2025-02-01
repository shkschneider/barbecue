package api

import (
	"barbecue/core"
	"barbecue/data"
)

func RemoveRecursive(task data.Task) error {
	children, err := GetChildren(task)
	if err != nil {
		core.Log.Debug("RemoveRecursive", err)
		return err
	} else if children != nil {
		for _, child := range *children {
			RemoveRecursive(child)
		}
	}
	Remove(task)
	return nil
}

func Remove(task data.Task) {
	core.Database.Delete(task)
}


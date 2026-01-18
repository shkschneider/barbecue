package api

import (
	"barbecue/core"
	"barbecue/data"
)

func Remove(task data.Task) error {
	children, err := GetChildren(task)
	if err != nil {
		core.Log.Debug("RemoveRecursive", err)
		return err
	} else if children != nil {
		for _, child := range *children {
			Remove(child)
		}
	}
	core.Database.Delete(task)
	return nil
}


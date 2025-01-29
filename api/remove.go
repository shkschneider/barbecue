package api

import (
	"barbecue/core"
	"barbecue/data"
)

func RemoveRecursive(task data.Task) {
	children, err := GetChildren(task)
	if err == nil {
		for _, child := range *children {
			RemoveRecursive(child)
		}
	}
	Remove(task)
}

func Remove(task data.Task) {
	core.Context.Database.Delete(task)
}


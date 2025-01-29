package server

import (
	"net/http"
	"strconv"
	"github.com/labstack/echo/v4"
	"barbecue/api"
	"barbecue/core"
	"barbecue/driver"
)

func Get(ctx echo.Context) error {
	out := driver.NewHtmlDriver()
	parents, err := api.GetParents()
	if err != nil {
		core.Context.Logger.Error(err)
		return out.Ko(ctx, http.StatusNotFound)
	}
	if parents == nil {
		return out.Ok(ctx, http.StatusOK, driver.T_INDEX, []ApiResponse{})
	}
	var responses []ApiResponse = *new([]ApiResponse)
	for _, parent := range *parents {
		children, _ := api.GetChildren(parent)
		responses = append(responses, ApiResponse {
			Parent: nil,
			Task: &parent,
			Children: children,
		})
	}
	return out.Ok(ctx, http.StatusOK, driver.T_INDEX, &responses)
}

func GetNew(ctx echo.Context) error {
	out := driver.NewHtmlDriver()
	return out.Ok(ctx, http.StatusOK, driver.T_FORM, nil)
}

func GetSlug(ctx echo.Context) error {
	out := driver.NewHtmlDriver()
	tasks, err := api.GetBySlug(ctx.Param("slug"))
	if err != nil {
		core.Context.Logger.Error(err)
		return out.Ko(ctx, http.StatusNotFound)
	}
	task := (*tasks)[0]
	core.Context.Logger.Debug(task)
	parent, _ := api.GetParent(task)
	children, _ := api.GetChildren(task)
	return out.Ok(ctx, http.StatusFound, driver.T_TASK, ApiResponse {
		Parent: parent,
		Task: &task,
		Children: children,
	})
}

func GetSlugEdit(ctx echo.Context) error {
	out := driver.NewHtmlDriver()
	tasks, err := api.GetBySlug(ctx.Param("slug"))
	if err != nil {
		core.Context.Logger.Error(err)
		return out.Ko(ctx, http.StatusNotFound)
	}
	task := (*tasks)[0]
	core.Context.Logger.Debug(task)
	return out.Ok(ctx, http.StatusFound, driver.T_FORM, ApiResponse {
		Task: &task,
	})
}

func GetSlugNew(ctx echo.Context) error {
	out := driver.NewHtmlDriver()
	parents, err := api.GetBySlug(ctx.Param("slug"))
	if err != nil {
		core.Context.Logger.Error(err)
		return out.Ko(ctx, http.StatusNotFound)
	}
	parent := (*parents)[0]
	core.Context.Logger.Debug(parent)
	return out.Ok(ctx, http.StatusOK, driver.T_FORM, ApiResponse {
		Parent: &parent,
	})
}

func GetSlugProgress(ctx echo.Context) error {
	out := driver.NewHtmlDriver()
	tasks, err := api.GetBySlug(ctx.Param("slug"))
	if err != nil {
		core.Context.Logger.Error(err)
		return out.Ko(ctx, http.StatusNotFound)
	}
	task := (*tasks)[0]
	if pc, err := strconv.ParseUint(ctx.Param("progress"), 10, 16) ; err == nil {
		if pc <= 0 {
			task.Progress = 0
		} else if pc >= 100 {
			task.Progress = 100
		} else {
			task.Progress = uint(pc)
		}
	}
	if err := core.Context.Database.Update(task) ; err != nil {
		return out.Ko(ctx, http.StatusInternalServerError)
	}
	core.Context.Logger.Debug(task)
	parent, err := api.GetParent(task)
	if err != nil {
		core.Context.Logger.Error(err)
		return out.Redirect(ctx, "/" + task.Slug)
	}
	return out.Redirect(ctx, "/" + parent.Slug)
}

func GetSlugDelete(ctx echo.Context) error {
	out := driver.NewHtmlDriver()
	tasks, err := api.GetBySlug(ctx.Param("slug"))
	if err != nil {
		core.Context.Logger.Error(err)
		return out.Ko(ctx, http.StatusNotFound)
	}
	task := (*tasks)[0]
	core.Context.Logger.Debug(task)
	api.RemoveRecursive(task)
	parent, err := api.GetParent(task)
	if err != nil {
		core.Context.Logger.Error(err)
		return out.Redirect(ctx, "/")
	}
	return out.Redirect(ctx, "/" + parent.Slug)
}

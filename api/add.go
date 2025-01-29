package api

import (
	"barbecue/core"
	"barbecue/data"
)

func Add(title string, description string) (*[]data.Task, error) {
	slug := data.Slugify(title)
	core.Context.Database.Insert(slug, title, description)
	tasks, err := GetBySlug(slug)
	return tasks, err
	// 	task, err := ctx.api.DB.add("", ctx.api.Request.Title, ctx.api.Request.Description)
	// 	if err != nil { return ctx.ko(http.StatusInternalServerError) }
	// 	return ctx.redirect("/" + task.Slug)
	// })
	// 	task, err := ctx.api.DB.add(ctx.api.Request.Slug, ctx.api.Request.Title, ctx.api.Request.Description)
	// 	if err != nil { return ctx.ko(http.StatusNotFound) }
	// 	return ctx.redirect("/" + task.Slug)
}

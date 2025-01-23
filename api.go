package main

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model			`json:"-"`
	Slug 		string 	`json:"slug",param:"slug"`
	Title 		string 	`json:"title",form:"title"`
	Description string 	`json:"description",form:"description"`
	Progress	uint	`json:"progress",form:"progress"`
	Super 		*uint	`json:"-"`
}

type (
	Api interface {
		Ok() error
		Ko() error
	}
	ApiContext struct {
		DB			*Database
		Request		ApiRequest
		Response	ApiResponse
	}
	ApiRequest struct {
		Id			uint 	`param:"id"`
		Slug		string 	`param:"slug"`
		Title 		string 	`form:"title"`
		Description string 	`form:"description"`
		Progress	uint 	`param:"progress"`
	}
	ApiResponse struct {
		Parent		*Task	`json:"parent,omitempty"`
		Task		*Task	`json:"task"`
		Children	*[]Task `json:"children"`
	}
	ApiError struct {
		Code		int
		Message		string
	}
)

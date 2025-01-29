package server

import (
	"barbecue/data"
)

type ApiRequest struct {
	Id			uint 	`param:"id"`
	Slug		string 	`param:"slug"`
	Title 		string 	`form:"title"`
	Description string 	`form:"description"`
	Progress	uint 	`param:"progress"`
}

type ApiResponse struct {
	Parent		*data.Task	`json:"parent,omitempty"`
	Task		*data.Task	`json:"task"`
	Children	*[]data.Task `json:"children"`
}

type ApiError struct {
	Code		int
	Message		string
}

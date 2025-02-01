package data

import (
	"strings"
	"gorm.io/gorm"
)

type Task struct {
  	gorm.Model				`json:"-"`
  	ID           	uint	`gorm:"primarykey;size:16"`
	Slug 			string 	`json:"slug",param:"slug"`
	Title 			string 	`json:"title",form:"title"`
	Description 	string 	`json:"description",form:"description"`
	Progress		uint	`json:"progress",form:"progress"`
	Super 			*uint	`json:"-"`
}

func NewTask() *Task {
	return &Task{}
}

func NewTasks(tasks ...Task) *[]Task {
	var data []Task
	for _, task := range tasks {
		data = append(data, task)
	}
	return &data
}

// "github.com/gosimple/slug"
func Slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.Map(func(r rune) rune {
		if (r < 'a' || r > 'z') && r != '-' && r != '_' && (r < '0' || r > '9') {
			return rune('-')
		} else {
			return r
		}
	}, s)
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, "--", "-")
	for strings.HasPrefix(s, "-") { s = strings.TrimPrefix(s, "-") }
	for strings.HasSuffix(s, "-") { s = strings.TrimSuffix(s, "-") }
	return s
}

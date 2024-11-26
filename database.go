package main

import (
	"fmt"
	// "github.com/gosimple/slug"
	"gorm.io/gorm"
	"os"
	"strings"
)

var db *gorm.DB

// Get

func slugify(s string) string {
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

func progress(task Task) uint {
	tasks, err := GetSubTasks(task.Slug)
	if err != nil || tasks == nil { return task.Progress }
	var p uint = 0
	for _, t := range *tasks {
		p += t.Progress
	}
	return (p / uint(len(*tasks)))
}

func GetParents() (*[]Task, error) {
	var tasks []Task
	result := db.Model(&Task{}).Where("parent_id IS NULL").Find(&tasks)
	if len(tasks) == 0 { return nil, nil }
	return &tasks, result.Error
}

func GetParent(slug string) (*Task, error) {
	task, err := GetTask(slug)
	if err != nil || task == nil { return nil, err }
	var parent Task
	db.Model(&Task{}).First(&parent, task.ParentID)
	// result := db.Model(&Task{}).First(&parent, task.ParentID)
	// if result.Error != nil || parent.ID == 0 { return nil, result.Error }
	return &parent, nil
}

func GetTask(slug string) (*Task, error) {
	var task Task
	result := db.Model(&Task{}).Where(Task { Slug: slug }).First(&task)
	return &task, result.Error
}

func GetSubTasks(slug string) (*[]Task, error) {
	task, err := GetTask(slug)
	if err != nil || task == nil { return nil, err }
	var tasks []Task
	result := db.Model(&Task{}).Where(Task { ParentID: &task.ID }).Order("progress").Find(&tasks)
	if len(tasks) == 0 { return nil, result.Error }
	return &tasks, result.Error
}

// Set

func New(slug string, title string, description string) (Task, error) {
	task := Task {
		Slug: slugify(title),
		Title: title,
		Description: description,
		ParentID: nil,
	}
	if len(slug) > 0 {
		if parent, err := GetTask(slug) ; err != nil {
			return Task{}, err
		} else {
			task.Slug = fmt.Sprintf("%v-%s", parent.ID, task.Slug)
			task.ParentID = &parent.ID
		}
	}
	result := db.Create(&task)
	return task, result.Error
}

func Update(task Task) (error) {
	result := db.Save(&task)
	return result.Error
}

func Delete(task Task) (error) {
	result := db.Delete(&task)
	return result.Error
}

func Dump() (*[]Task) {
	if os.Getenv("DEBUG") != "true" {
		return nil
	}
	var tasks []Task
	db.Model(&Task{}).Find(&tasks)
	return &tasks
}

// Main

func Database(d gorm.Dialector) (*gorm.DB, error) {
	var err error
	db, err = gorm.Open(d, &gorm.Config{}) ; if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Task{})
	if os.Getenv("DEBUG") == "true" {
		var task Task
		var tasks []Task
		db.Session(&gorm.Session { AllowGlobalUpdate: true }).Delete(&Task{})
		New("", "Job", "You can't get a job without experience.")
		New("job", "Experience", "You can't get experience without a job.")
		New("", "Ask", "Ask a question.")
		task, _ = New("ask", "Google", "First answer is 'Google it'!")
		New(task.Slug, "Search Results", "First link on Google is the exact page where you asked your initial question...")
		New("", "Definition", "Definition of recursion:")
		task, _ = New("definition", "Recursion", "See its definition...")
		task.Description = "See its [definition](/definition)..."
		Update(task)
		db.Model(&Task{}).Find(&tasks)
		for _, task := range tasks {
			fmt.Println(task.ID, task.Slug)
		}
	}
	return db, nil
}

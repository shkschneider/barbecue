package main

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
)

type Database struct {
	*gorm.DB
}

// "github.com/gosimple/slug"
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

// Get

func (d Database) getParents() (*[]Task, error) {
	var parents []Task
	result := d.Model(&Task{}).Where("super IS NULL").Order("progress").Find(&parents)
	if len(parents) == 0 { return nil, nil }
	return &parents, result.Error
}

func (d Database) get(slug string) (*Task, error) {
	var task Task
	result := d.Model(&Task{}).Where(Task { Slug: slug }).First(&task)
	if result.Error == nil {
		var children []Task
		result := d.Model(&Task{}).Where(Task { Super: &task.ID }).Find(&children)
		if result.Error == nil && len(children) > 0 {
			task.Progress = 0
			for _, child := range children {
				task.Progress += child.Progress
			}
			task.Progress = task.Progress / uint(len(children))
			d.Save(&task)
		}
	}
	return &task, result.Error
}

func (d Database) getAll() (*[]Task, error) {
	var tasks []Task
	result := d.Model(&Task{}).Find(&tasks)
	if len(tasks) == 0 { return nil, nil }
	return &tasks, result.Error
}

func (d Database) getParent(task *Task) (*Task, error) {
	var parent Task
	result := d.Model(&Task{}).First(&parent, task.Super)
	if result.Error != nil || parent.ID == 0 { return nil, result.Error }
	return &parent, nil
}

func (d Database) getChildren(task *Task) (*[]Task, error) {
	var children []Task
	result := d.Model(&Task{}).Where(Task { Super: &task.ID }).Order("progress").Find(&children)
	if len(children) == 0 { return nil, result.Error }
	return &children, result.Error
}

// Set

func (d *Database) add(slug string, title string, description string) (Task, error) {
	task := Task {
		Slug: slugify(title),
		Title: title,
		Description: description,
		Super: nil,
	}
	if len(slug) > 0 {
		if parent, err := d.get(slug) ; err != nil {
			return Task{}, err
		} else {
			task.Slug = fmt.Sprintf("%v-%s", parent.ID, task.Slug)
			task.Super = &parent.ID
		}
	}
	result := d.Create(&task)
	return task, result.Error
}

func (d *Database) update(task Task) (error) {
	result := d.Save(&task)
	return result.Error
}

func (d *Database) remove(task Task) (error) {
	result := d.Delete(&task)
	return result.Error
}

// Main

func NewDatabase(d gorm.Dialector) (*Database, error) {
	var db Database
	config := gorm.Config {
		Logger: logger.Default.LogMode(logger.Silent),
	}
	if d, err := gorm.Open(d, &config) ; err != nil {
		return nil, err
	} else {
		db = Database { d }
	}
	db.AutoMigrate(&Task{})
	if DEBUG {
		var task Task
		var tasks []Task
		db.Session(&gorm.Session { AllowGlobalUpdate: true }).Delete(&Task{})
		db.add("", "Job", "You can't get a job without experience.")
		db.add("job", "Experience", "You can't get experience without a job.")
		db.add("", "Ask", "Ask a question.")
		task, _ = db.add("ask", "Google", "First answer is 'Google it'!")
		db.add(task.Slug, "Search Results", "First link on Google is the exact page where you asked your initial question...")
		db.add("", "Definition", "Definition of recursion:")
		task, _ = db.add("definition", "Recursion", "See its definition...")
		task.Description = "See its [definition](/definition)..."
		db.update(task)
		db.Model(&Task{}).Find(&tasks)
		for _, task := range tasks {
			log.Debug(fmt.Sprintf("#%d %s", task.ID, task.Slug))
		}
	}
	return &db, nil
}

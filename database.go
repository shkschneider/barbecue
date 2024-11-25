package main

import (
	"fmt"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
	"os"
)

var db *gorm.DB

// Get

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

func GetTask(slug string) (*Task, error) {
	var task Task
	result := db.Model(&Task{}).Where(Task { Slug: slug }).First(&task)
	return &task, result.Error
}

func GetParent(slug string) (*Task, error) {
	task, err := GetTask(slug)
	if err != nil || task == nil { return nil, err }
	var parent Task
	result := db.Model(&Task{}).First(&parent, task.ParentID)
	if parent.ID == 0 { return nil, result.Error }
	return &parent, result.Error
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

// Main

func Database(d gorm.Dialector) (*gorm.DB, error) {
	var err error
	db, err = gorm.Open(d, &gorm.Config{}) ; if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Task{})
	if os.Getenv("DEBUG") == "true" {
		fmt.Println(slugify("This #Is_A_Slugify Test!!!"))
		fmt.Println(slug.Make("This #Is_A_Slug tESt!!!"))
		db.Session(&gorm.Session { AllowGlobalUpdate: true }).Delete(&Task{})
		New("", "Title1", "Description _One_")
		New("", "Title2", "Description _Two_")
		New("", "Title3", "Description _Three_")
		New("title1", "Title11", "Description11")
		New("title1", "Title12", "Description12")
		New("title1", "Title13", "Description13")
		var tasks []Task
		db.Model(&Task{}).Find(&tasks)
		for _, task := range tasks {
			fmt.Println(task.ID, task.Slug)
		}
	}
	return db, nil
}

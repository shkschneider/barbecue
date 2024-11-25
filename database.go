package main

import (
	"fmt"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func DatabaseGet(slug string) (*Response, error) {
	var task Task
	result := DB.Model(&Task{}).Where(Task { Slug: slug }).First(&task)
	var parent *Task = new(Task)
	DB.Model(&Task{}).First(parent, task.ParentID)
	if parent.ID == 0 { parent = nil }
	var tasks *[]Task = new([]Task)
	DB.Model(&Task{}).Where(Task { ParentID: &task.ID }).Order("progress").Find(tasks)
	if len(*tasks) == 0 {
		tasks = nil
	} else {
		var p uint = 0 ; for _, t := range *tasks { p += t.Progress }
		task.Progress = (p / uint(len(*tasks)))
		DB.Save(&task)
	}
	return &Response {
		Parent: parent,
		Task: &task,
		SubTasks: tasks,
	}, result.Error
}

func DatabaseSet(slug string, title string, description string) (Task, error) {
	var response *Response
	if len(slug) > 0 {
		if r, err := DatabaseGet(slug) ; err != nil {
			return Task{}, err
		} else {
			response = r
		}
	} else {
		response = &Response{}
	}
	task := Task {
		Slug: slugify(title),
		Title: title,
		Description: description,
		ParentID: nil,
	}
	if response.Task != nil {
		task.Slug = fmt.Sprintf("%v-%s", response.Task.ID, task.Slug)
		task.ParentID = &response.Task.ID
	}
	result := DB.Create(&task)
	return task, result.Error
}

func Database(d gorm.Dialector) (*gorm.DB, error) {
	db, err := gorm.Open(d, &gorm.Config{}) ; if err != nil {
		return nil, err
	} else {
		DB = db
		DB.AutoMigrate(&Task{})
	}
	if os.Getenv("DEBUG") == "true" {
		fmt.Println(slugify("This #Is_A_Slugify Test!!!"))
		fmt.Println(slug.Make("This #Is_A_Slugi Test!!!"))
		DB.Session(&gorm.Session { AllowGlobalUpdate: true }).Delete(&Task{})
		DatabaseSet("", "Title1", "Description _One_")
		DatabaseSet("", "Title2", "Description _Two_")
		DatabaseSet("", "Title3", "Description _Three_")
		DatabaseSet("title1", "Title11", "Description11")
		DatabaseSet("title1", "Title12", "Description12")
		DatabaseSet("title1", "Title13", "Description13")
		var tasks []Task
		DB.Model(&Task{}).Find(&tasks)
		for _, task := range tasks {
			fmt.Println(task.ID, task.Slug)
		}
	}
	return DB, nil
}

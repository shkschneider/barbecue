package data

import (
	"os"
	"strings"
	"path/filepath"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/driver/sqlite"
)

type Database struct {
	Path	string
	Orm		*gorm.DB
}

func NewDatabase(path string, name string) (*Database, error) {
	// path
	home, err := os.UserHomeDir()
	if err != nil {
		path = name + ".sqlite"
	} else if len(path) == 0 {
		path = name + ".sqlite"
		if err == nil {
	    	path = filepath.Join(home, ".local", "state", name + ".sqlite")
		} else {
			path = name + ".sqlite"
		}
	} else {
		path = strings.Replace(path, "~", home, 1)
	}
	// orm
	driver := sqlite.Open(path)
	orm, err := gorm.Open(driver, &gorm.Config {
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	// db
	var db Database = Database {
		Path: path,
		Orm: orm,
	}
	orm.AutoMigrate(&Task{})
	if os.Getenv("DEBUG") == "true" {
		var task Task
		var tasks []Task
		orm.Session(&gorm.Session { AllowGlobalUpdate: true }).Delete(&Task{})
		db.Insert("", "Job", "You can't get a job without experience.")
		db.Insert("job", "Experience", "You can't get experience without a job.")
		db.Insert("", "Ask", "Ask a question.")
		task, _ = db.Insert("ask", "Google", "First answer is 'Google it'!")
		db.Insert(task.Slug, "Search Results", "First link on Google is the exact page where you asked your initial question...")
		db.Insert("", "Definition", "Definition of recursion:")
		task, _ = db.Insert("definition", "Recursion", "See its definition...")
		task.Description = "See its [definition](/definition)..."
		db.Update(task)
		orm.Model(&Task{}).Find(&tasks)
	}
	return &db, nil
}

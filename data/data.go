package data

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	// "barbecue/core"
)

type Database struct {
	Orm	*gorm.DB
}

func New(d gorm.Dialector, debug bool) (*Database, error) {
	config := gorm.Config {
		Logger: logger.Default.LogMode(logger.Silent),
	}
	orm, err := gorm.Open(d, &config)
	if err != nil {
		return nil, err
	}
	var db Database = Database {
		Orm: orm,
	}
	orm.AutoMigrate(&Task{})
	if debug == true {
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
		// for _, task := range tasks {
		// 	log.Debug(fmt.Sprintf("#%d %s", task.ID, task.Slug))
		// }
	}
	return &db, nil
}

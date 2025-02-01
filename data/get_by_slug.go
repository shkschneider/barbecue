package data

func (db *Database) GetBySlug(slug string) (*[]Task, error) {
	var tasks []Task
	if result := db.Orm.Model(&Task{}).Where(Task { Slug: slug }).Find(&tasks) ;
		result.Error != nil || len(tasks) == 0 {
			return nil, result.Error
	}
	task := tasks[0]
	var children []Task
	if result := db.Orm.Model(&Task{}).Where(Task { Super: &task.ID }).Find(&children) ;
		result.Error != nil {
			return &tasks, result.Error
	}
	if len(children) > 0 {
		task.Progress = 0
		for _, child := range children {
			task.Progress += child.Progress
		}
		task.Progress = task.Progress / uint(len(children))
		db.Update(task)
	}
	return &tasks, nil
}

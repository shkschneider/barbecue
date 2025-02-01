package data

func (db *Database) GetById(id uint) (*Task, error) {
	var task Task
	if result := db.Orm.Model(&Task{}).Where(struct { ID uint } { ID: id }).First(&task) ;
		result.Error != nil || task.ID == 0 {
			return nil, result.Error
	}
	var children []Task
	if result := db.Orm.Model(&Task{}).Where(Task { Super: &task.ID }).Find(&children) ;
		result.Error != nil {
			return &task, result.Error
	}
	if len(children) > 0 {
		task.Progress = 0
		for _, child := range children {
			task.Progress += child.Progress
		}
		task.Progress = task.Progress / uint(len(children))
		db.Update(task)
	}
	return &task, nil
}

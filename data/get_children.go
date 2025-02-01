package data

func (db *Database) GetChildren(task Task) (*[]Task, error) {
	var children []Task
	if result := db.Orm.Model(&Task{}).Where(Task { Super: &task.ID }).Order("progress").Find(&children) ;
		result.Error != nil || len(children) == 0 {
			return nil, result.Error
	}
	return &children, nil
}

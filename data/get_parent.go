package data

func (db *Database) GetParent(task Task) (*Task, error) {
	if task.Super == nil { return nil, nil }
	var parent Task
	if result := db.Orm.Model(&Task{}).First(&parent, task.Super) ;
		result.Error != nil || parent.ID == 0 {
			return nil, result.Error
	}
	return &parent, nil
}

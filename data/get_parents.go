package data

func (db *Database) GetParents() (*[]Task, error) {
	var parents []Task
	result := db.Orm.Model(&Task{}).Where("super IS NULL").Order("progress DESC").Find(&parents)
	return &parents, result.Error
}

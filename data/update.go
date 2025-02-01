package data

func (db *Database) Update(task Task) error {
	result := db.Orm.Save(&task)
	return result.Error
}

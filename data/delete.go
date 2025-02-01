package data

func (db *Database) Delete(task Task) {
	db.Orm.Delete(&task)
}

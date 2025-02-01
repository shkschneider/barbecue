package data

func (db *Database) Insert(slug string, title string, description string) (Task, error) {
	task := Task {
		Slug: slug,
		Title: title,
		Description: description,
		Super: nil,
	}
	if len(slug) > 0 {
		if parents, err := db.GetBySlug(slug) ; err != nil {
			return Task{}, err
		} else if parents != nil && len(*parents) > 0 {
			parent := (*parents)[0]
			//task.Slug = fmt.Sprintf("%v-%s", parent.ID, task.Slug)
			task.Super = &parent.ID
		}
	}
	result := db.Orm.Create(&task)
	return task, result.Error
}

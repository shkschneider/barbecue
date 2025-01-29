package data

func (db *Database) GetParents() (*[]Task, error) {
	var parents []Task
	result := db.Orm.Model(&Task{}).Where("super IS NULL").Order("progress").Find(&parents)
	return &parents, result.Error
}

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

func (db *Database) GetBySlug(slug string) (*[]Task, error) {
	var tasks []Task
	if result := db.Orm.Model(&Task{}).Where("Slug LIKE ?", "%" + slug + "%").Find(&tasks) ;
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

func (db *Database) GetAll() (*[]Task, error) {
	var tasks []Task
	if result := db.Orm.Model(&Task{}).Find(&tasks) ;
		result.Error != nil {
			return nil, result.Error
	}
	return &tasks, nil
}

func (db *Database) GetParent(task Task) (*Task, error) {
	var parent Task
	if result := db.Orm.Model(&Task{}).First(&parent, task.Super) ;
		result.Error != nil || parent.ID == 0 {
			return nil, result.Error
	}
	return &parent, nil
}

func (db *Database) GetChildren(task Task) (*[]Task, error) {
	var children []Task
	if result := db.Orm.Model(&Task{}).Where(Task { Super: &task.ID }).Order("progress").Find(&children) ;
		result.Error != nil || len(children) == 0 {
			return nil, result.Error
	}
	return &children, nil
}

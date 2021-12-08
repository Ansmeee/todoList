package todo

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"todoList/src/models/todo"
	"todoList/src/services/common"
	"todoList/src/utils/database"
)

type TodoService struct{}

var ctx = context.Background()
var Model = &todo.TodoModel{}

func (TodoService) NewModel() *todo.TodoModel {
	return new(todo.TodoModel)
}

func (TodoService) Create(data *todo.TodoModel) (todo *todo.TodoModel, error error) {
	db := database.Connect("")
	defer database.Close(db)

	uid, error := common.GetUID()
	if error != nil {
		return
	}

	data.Id = uid
	error = db.Model(Model).Create(data).Error
	if error != nil {
		return
	}

	todo = data
	return
}

func (TodoService) FindByID(id string) (data todo.TodoModel, error error) {
	db := database.Connect("")
	defer database.Close(db)

	error = db.Model(Model).First(&data, id).Error
	if error != nil {
		if errors.Is(error, gorm.ErrRecordNotFound) {
			error = errors.New("该记录不存在")
		}

		return
	}

	return
}

type QueryForm struct {
	Keywords string `json:"keywords" form:"keywords"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
}

func (TodoService) List(form *QueryForm) (data []todo.TodoModel, total int64, error error) {
	db := database.Connect("")
	defer database.Close(db)

	db = db.Model(Model)
	if len(form.Keywords) > 0 {
		db = db.Where("title like ?", "%" + form.Keywords + "%")
	}

	db.Count(&total)
	if total == 0 {
		return
	}

	limit := form.PageSize
	offset := (form.Page - 1) * limit
	error = db.Limit(limit).Offset(offset).Find(&data).Error
	if error != nil {
		return
	}

	return
}

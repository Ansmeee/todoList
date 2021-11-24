package todo

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"todoList/src/models/todo"
	"todoList/src/utils/database"
)

type TodoService struct {}
var ctx = context.Background()
var Model = &todo.TodoModel{}

func (TodoService) NewModel() (*todo.TodoModel)  {
	return new(todo.TodoModel)
}

func (TodoService) Create(data *todo.TodoModel) (todo *todo.TodoModel, error error) {
	db := database.Connect("")
	defer database.Close(db)

	error = db.Model(Model).Create(data).Error
	if error != nil {
		return
	}

	todo = data
	return
}

func (TodoService) FindByID(id uint) (data todo.TodoModel, error error)  {
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

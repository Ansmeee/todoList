package todo

import (
	"context"
	"todoList/src/models/todo"
	"todoList/src/services/common"
	"todoList/src/utils/database"
)

type TodoService struct{}

var ctx = context.Background()
var model = &todo.TodoModel{}
var service = &TodoService{}

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
	error = db.Model(model).Create(data).Error
	if error != nil {
		return
	}

	todo = data
	return
}

type UpdateForm struct {
	Id      string `form:"id"`
	Title   string `form:"title"`
	Content string `form:"content"`
}
func (TodoService) Update(todo *todo.TodoModel, data *UpdateForm) (error error) {
	db := database.Connect("")
	defer database.Close(db)

	updateData := map[string]interface{}{
		"title":   data.Title,
		"content": data.Content,
	}

	error = db.Model(todo).Where("uid = ?", data.Id).Updates(updateData).Error
	return
}

func (TodoService) UpdateAttr(todo *todo.TodoModel, attrName, attrValue string) (error error)  {
	db := database.Connect("")
	defer database.Close(db)

	error = db.Model(todo).Where("uid = ?", todo.Id).Update(attrName, attrValue).Error
	return
}

func (TodoService) FindByID(id string) (todo *todo.TodoModel, error error) {
	db := database.Connect("")
	defer database.Close(db)

	todo = service.NewModel()
	error = db.Where("uid = ?", id).Find(todo).Error
	if error != nil {
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

	db = db.Model(model)
	if len(form.Keywords) > 0 {
		db = db.Where("title like ?", "%"+form.Keywords+"%")
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

func (TodoService) Delete(todo *todo.TodoModel) (error error) {
	db := database.Connect("")
	defer database.Close(db)

	error = db.Where("uid = ?", todo.Id).Delete(todo).Error
	if error != nil {
		return
	}

	return
}

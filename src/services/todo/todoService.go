package todo

import (
	"context"
	"fmt"
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

	uid := common.GetUID("todoUID")

	data.Id = uid
	error = db.Model(model).Create(data).Error
	if error != nil {
		return
	}

	todo = data
	return
}

type UpdateForm struct {
	Id       int `form:"id"`
	Title    string `form:"title"`
	Content  string `form:"content"`
	Priority int    `form:"priority"`
	Deadline string `form:"deadline"`
}

func (TodoService) Update(todo *todo.TodoModel, data *UpdateForm) (error error) {
	db := database.Connect("")
	defer database.Close(db)

	updateData := map[string]interface{}{
		"title":    data.Title,
		"content":  data.Content,
		"priority": data.Priority,
		"deadline": data.Deadline,
	}

	error = db.Model(todo).Where("uid = ?", data.Id).Updates(updateData).Error
	return
}

func (TodoService) UpdateAttr(todo *todo.TodoModel, attrName string, attrValue interface{}) (error error) {
	db := database.Connect("")
	defer database.Close(db)

	if attrName == "top" && attrValue == "1" {
		attrValue = common.Incr("top")
	}

	error = db.Model(todo).Where("uid = ?", todo.Id).Update(attrName, attrValue).Error
	return
}

func (TodoService) FindByID(id int) (todo *todo.TodoModel, error error) {
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
	From      string   `json:"from" form:"from"`
	Keywords  string   `json:"keywords" form:"keywords"`
	Page      int      `json:"page" form:"page"`
	PageSize  int      `json:"page_size" form:"page_size"`
	SortBy    string   `json:"sort_by" form:"sort_by"`
	SortOrder string   `json:"sort_order" form:"sort_order"`
	Rules     []string `json:"rules" form:"rules"`
	Wheres    [][]string
	ListId    string
}

func (TodoService) List(form *QueryForm) (data []todo.TodoModel, total int64, error error) {
	db := database.Connect("")
	defer database.Close(db)

	db = db.Model(model)
	if len(form.ListId) > 0 {
		db = db.Where("list_id = ?", form.ListId)
	}

	if len(form.Wheres) > 0 {
		for _, where := range form.Wheres {
			if where[1] == "=" {
				db = db.Where(map[string]interface{}{where[0]: where[2]})
			}

			if where[1] == "<>" {
				db = db.Not(map[string]interface{}{where[0]: where[2]})
			}

		}
	}

	if len(form.Keywords) > 0 {
		db = db.Where("title like ?", "%"+form.Keywords+"%")
	}

	db.Count(&total)
	if total == 0 {
		return
	}

	if form.SortBy != "" && form.SortOrder != "" {
		db = db.Order(fmt.Sprintf("`%s` %s", form.SortBy, form.SortOrder))
	}

	if form.PageSize > 0 {
		limit := form.PageSize
		offset := (form.Page - 1) * limit
		db = db.Limit(limit).Offset(offset)
	}

	error = db.Find(&data).Error
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

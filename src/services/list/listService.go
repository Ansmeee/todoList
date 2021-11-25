package list

import (
	"todoList/src/models/list"
	"todoList/src/utils/database"
)

type ListService struct {}
var model = &list.ListModel{}

func (ListService) NewModel() *list.ListModel {
	return new(list.ListModel)
}

func (ListService) Create(list *list.ListModel) (data *list.ListModel, error error) {
	db := database.Connect("")
	defer database.Close(db)

	error = db.Model(model).Create(list).Error
	if error != nil {
		return
	}

	data = list
	return
}
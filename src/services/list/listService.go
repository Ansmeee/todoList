package list

import (
	"errors"
	"todoList/src/models/list"
	"todoList/src/utils/database"
)

type ListService struct {}
var model = &list.ListModel{}
var service = &ListService{}

func (ListService) NewModel() *list.ListModel {
	return new(list.ListModel)
}

func (ListService) FindByID(id uint) (list *list.ListModel, error error) {
	db := database.Connect("")
	defer database.Close(db)

	list = service.NewModel()
	error = db.Model(model).Find(list, id).Error
	if error != nil {
		return
	}

	return
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

func (ListService) Update(list, data *list.ListModel) (result *list.ListModel, error error) {
	db := database.Connect("")
	defer database.Close(db)

	if list.Id == 0 {
		error = errors.New("清单不存在")
		return
	}

	error = db.Model(list).Updates(data).Error
	if error != nil {
		return
	}

	result = list
	return
}

func (ListService) Delete(list *list.ListModel) (error error) {
	db := database.Connect("")
	defer database.Close(db)

	if list.Id == 0 {
		error = errors.New("清单不存在")
		return
	}

	error = db.Model(model).Delete(list).Error
	return
}

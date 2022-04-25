package list

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"todoList/src/models/list"
	"todoList/src/models/todo"
	"todoList/src/services/common"
	"todoList/src/utils/database"
)

type ListService struct{}

var model = &list.ListModel{}
var service = &ListService{}

func (ListService) NewModel() *list.ListModel {
	return new(list.ListModel)
}

func (ListService) FindByID(id string) (list *list.ListModel, error error) {
	db := database.Connect("")
	defer database.Close(db)

	list = service.NewModel()
	error = db.Model(model).Where("uid = ?", id).Find(list).Error
	if error != nil {
		return
	}

	return
}

type Params struct {
	Keywords string
	PageSize int
	Page     int
	Userid   string
}

func (ListService) List(params *Params) (total int64, data []*list.ListModel, error error) {
	db := database.Connect("")
	defer database.Close(db)

	total = 0
	data = []*list.ListModel{}
	page := (params.Page - 1) * params.PageSize

	query := db.Model(model)

	if params.Userid != "" {
		query = query.Where("user_id = ?", params.Userid)
	}

	if len(params.Keywords) > 0 {
		query = query.Where("title like ?", "%"+params.Keywords+"%")
	}

	if query.Count(&total).Error != nil {
		error = errors.New("获取失败")
		return
	}

	if total == 0 {
		return
	}

	if query.Limit(params.PageSize).Offset(page).Find(&data).Error != nil {
		error = errors.New("获取失败")
		return
	}

	return
}

func (ListService) Create(list *list.ListModel) (data *list.ListModel, error error) {
	db := database.Connect("")
	defer database.Close(db)

	list.Id = common.GetUID()
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

	if len(list.Id) == 0 {
		error = errors.New("清单不存在")
		return
	}

	updateData := map[string]interface{}{"updated_at": time.Now().Format("2006-01-02 15:01:05")}
	if data.Title != list.Title {
		updateData["title"] = data.Title
	}

	error = db.Model(list).Where("uid = ?", data.Id).Updates(updateData).Error

	if error != nil {
		return
	}

	result = list
	return
}

func (ListService) Delete(list *list.ListModel) error {
	db := database.Connect("")
	defer database.Close(db)

	if len(list.Id) == 0 {
		return errors.New("清单不存在")
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		updateData := map[string]interface{}{
			"list_id": 0,
			"type":    "",
		}
		if err := db.Model(&todo.TodoModel{}).Where("list_id = ?", list.Id).Updates(updateData).Error; err != nil {
			return err
		}

		if err := db.Where("uid = ?", list.Id).Delete(list).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

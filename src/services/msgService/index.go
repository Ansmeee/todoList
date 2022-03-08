package msgService

import (
	"todoList/src/models/msgModel"
	"todoList/src/models/user"
	"todoList/src/utils/database"
)

type MsgService struct{}


func (MsgService) NewMsgModel() (*msgModel.MsgModel) {
	return new(msgModel.MsgModel)
}

func (MsgService) FindByID(id string) (msg *msgModel.MsgModel, error error) {
	db := database.Connect("")
	defer database.Close(db)

	msg = MsgService{}.NewMsgModel()
	error = db.Where("uid = ?", id).Find(msg).Error
	if error != nil {
		return
	}

	return
}

func (MsgService) UnreadCount() int64 {
	db := database.Connect("")
	defer database.Close(db)

	var count int64 = 0
	user := user.User()
	db.Model(MsgService{}.NewMsgModel()).Where("user_id = ? and status = ?", user.Id, msgModel.StatusUnread).Count(&count)
	return count
}

type ListForm struct {
	PageSize int `form:"page_size"`
	Page     int `form:"page"`
}

func (MsgService) List(form *ListForm) (data []msgModel.MsgModel, error error) {
	db := database.Connect("")
	defer database.Close(db)

	db = db.Model(MsgService{}.NewMsgModel())
	db = db.Order("`status`").Order("`id` desc")
	page, pageSize := paginate(form.Page, form.PageSize)
	error = db.Limit(pageSize).Offset(page).Find(&data).Error
	return
}

type AttrForm struct {
	Id    string  `form:"id"`
	Name  string `form:"name"`
	Value string `form:"value"`
}
func (MsgService) Update(msg *msgModel.MsgModel, attrName string, attrValue interface{}) (error error) {
	db := database.Connect("")
	defer database.Close(db)

	updateData := map[string]interface{}{attrName: attrValue}
	error = db.Model(msg).Where("uid = ?", msg.Id).Updates(updateData).Error
	return
}

func paginate(formPage, formPageSize int) (int, int) {
	page := 1
	if formPage > 0 {
		page = formPage
	}

	pageSize := 20
	if formPageSize > 0 {
		pageSize = formPageSize
	}

	page = (page - 1) * pageSize
	return page, pageSize
}
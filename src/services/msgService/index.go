package msgService

import (
	"errors"
	"fmt"
	"todoList/src/models/msgModel"
	"todoList/src/models/user"
	"todoList/src/utils/database"
)

type MsgService struct{}

var service = &MsgService{}

func (MsgService) NewMsgModel() *msgModel.MsgModel {
	return new(msgModel.MsgModel)
}

func (MsgService) FindByID(id string) (msg *msgModel.MsgModel, error error) {
	db := database.Connect("")
	defer database.Close(db)

	msg = service.NewMsgModel()
	error = db.Where("uid = ?", id).Find(msg).Limit(1).Error
	if error != nil {
		fmt.Println(error.Error())
		return
	}

	return
}

func (MsgService) UnreadCount() int64 {
	db := database.Connect("")
	defer database.Close(db)

	var count int64 = 0
	user := user.User()
	db.Model(MsgService{}.NewMsgModel()).Where("user_id = ? and status = ?", user.Id, msgModel.STATUS_UNREAD).Count(&count)
	return count
}

type ListForm struct {
	PageSize int `form:"page_size"`
	Page     int `form:"page"`
	Status   int `form:"status"`
	Force    int `form:"force"`
}

func (MsgService) List(form *ListForm) (data []msgModel.MsgModel, error error) {
	db := database.Connect("")
	defer database.Close(db)

	userID := user.User().Id
	db = db.Model(MsgService{}.NewMsgModel()).Where("user_id = ?", userID)

	avaSMap := map[int]bool{2: true, 1: true}

	if _, ok := avaSMap[form.Force]; ok {
		fmt.Println("force")
		db = db.Where("`force` = ?", form.Force)
	}

	avaFMap := map[int]bool{2: true, 1: true}
	if _, ok := avaFMap[form.Status]; ok {
		fmt.Println("status")
		db = db.Where("`status` = ?", form.Status)
	}

	db = db.Order("`status`").Order("`id` desc")
	page, pageSize := paginate(form.Page, form.PageSize)
	error = db.Limit(pageSize).Offset(page).Find(&data).Error
	return
}

type AttrForm struct {
	Id    string `form:"id"`
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

func (MsgService) Create(data *msgModel.MsgModel) error {
	db := database.Connect("")
	defer database.Close(db)

	if error := db.Model(&msgModel.MsgModel{}).Create(data).Error; error != nil {
		fmt.Println("MsgService Create Error:", error.Error())
		return errors.New("????????????")
	}

	return nil
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

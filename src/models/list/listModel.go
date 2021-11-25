package list

import "todoList/src/models"

type ListModel struct {
	Title    string `json:"title" form:"title"`
	Label    string `json:"label" form:"label"`
	CreateId uint   `json:"create_id" form:"create_id"`
	models.Model
}

func (ListModel) TableName() string {
	return "list"
}

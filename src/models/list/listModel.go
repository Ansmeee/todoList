package list

import "todoList/src/models"

type ListModel struct {
	Title    string `json:"title" form:"title"`
	Color    string `json:"color" form:"color"`
	Hide     uint `json:"hide" form:"hide"`
	Type     string `json:"type" form:"type"`
	CreateId uint   `json:"create_id" form:"create_id"`
	models.Model
}

func (ListModel) TableName() string {
	return "list"
}

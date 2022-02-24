package list

import "todoList/src/models"

type ListModel struct {
	Title  string `json:"title" form:"title"`
	Color  string `json:"color" form:"color"`
	Hide   uint   `json:"hide" form:"hide"`
	Type   string `json:"type" form:"type"`
	UserId int64    `json:"user_id" gorm:"<-:create"`
	models.Model
}

func (ListModel) TableName() string {
	return "list"
}

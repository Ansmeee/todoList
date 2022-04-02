package msgModel

import "todoList/src/models"

type MsgModel struct {
	UserId  string `json:"user_id"`
	Content string `json:"content"`
	Status  int    `json:"status"`
	Force   int    `json:"force"`
	Link    string `json:"link"`
	models.Model
}

const (
	StatusUnread = 1
	StatusRead   = 2

	Force = 2
)

func (MsgModel) TableName() string {
	return "user_msg"
}

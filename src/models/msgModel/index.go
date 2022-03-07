package msgModel

import "todoList/src/models"

type MsgModel struct {
	UserId  string `json:"user_id"`
	Content string `json:"content"`
	Status  int    `json:"status"`
	models.Model
}

const (
	StatusUnread = 0
	StatusRead   = 1
)

func (MsgModel) TableName() string {
	return "user_msg"
}
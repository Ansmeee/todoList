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
	STATUS_UNREAD = 1
	STATUS_READ   = 2

	FORCE    = 2
	UN_FORCE = 1
)

func (MsgModel) TableName() string {
	return "user_msg"
}

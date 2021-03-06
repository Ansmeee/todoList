package todo

import "todoList/src/models"

type TodoModel struct {
	Title    string `json:"title" form:"title"`
	Type     string `json:"type" form:"type"`
	Content  string `json:"content" form:"content"'`
	Status   int    `json:"status" form:"status"`
	ParentId int    `json:"parent_id" form:"parent_id"`
	ListId   string `json:"list_id" form:"list_id"`
	UserId   string `json:"user_id" form:"user_id"`
	Tags     string `json:"tags" form:"tags"`
	Priority int    `json:"priority" form:"priority"`
	Top      int    `json:"top" form:"top"`
	Deadline string `json:"deadline" form:"deadline"`
	models.Model
}

const (
	STATUS_ACTIVE = 1
	STATUS_DONE   = 2
)

func (TodoModel) TableName() string {
	return "todo"
}

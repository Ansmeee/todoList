package todo

import "todoList/src/models"

type TodoModel struct {
	Title        string `json:"title" form:"title"`
	Type         string `json:"type" form:"type"`
	UserId       int    `json:"user_id" form:"user_id"`
	Deadline     string `json:"deadline" form:"deadline"`
	Status       string `json:"status" form:"status"`
	Remark       string `json:"remark" form:"remark"`
	models.Model `form:"models_model"`
}

func (TodoModel) TableName() string {
	return "todo"
}

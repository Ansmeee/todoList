package user

import "todoList/src/models"

type UserModel struct {
	Name  string `json:"name" form:"name"`
	Email string `json:"email" form:"email" gorm:"unique"`
	Phone string `json:"phone" form:"phone"`
	Icon  string `json:"icon" form::"icon"`
	models.Model
}

func (UserModel) TableName() string {
	return "user"
}

func (UserModel) Active() bool {
	return true
}

package user

import (
	"todoList/src/models"
)

type AuthModel struct {
	Email string `json:"email" form:"email" gorm:"unique"`
	Auth  string `json:"auth" form:"auth"`
	models.Model
}

var user *UserModel

func (AuthModel) TableName() string {
	return "user_auth"
}

func (AuthModel) Login(userInfo *UserModel) {
	user = userInfo
}

func (AuthModel) GetUser() *UserModel {
	return user
}
package user

import "todoList/src/models"

type UserModel struct {
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email" gorm:"unique"`
	Phone    string `json:"phone" form:"phone"`
	Icon     string `json:"icon" form:"icon"`
	Verified int    `json:"verified" form:"verify"`
	models.Model
}

const (
	EMAIL_VERIFIED = 2
	PHONE_VERIFIED = 1
)

func (UserModel) TableName() string {
	return "user"
}

func (UserModel) Active() bool {
	return true
}

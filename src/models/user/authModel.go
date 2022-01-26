package user

import (
	"time"
)

type AuthModel struct {
	Account   int       `json:"account" form:"account" gorm:"unique"`
	Auth      string    `json:"auth" form:"auth"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;<-:create"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
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

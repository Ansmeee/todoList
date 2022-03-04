package user

import (
	"time"
)

type AuthModel struct {
	Account   string    `json:"account" form:"account" gorm:"unique"`
	Auth      string    `json:"auth" form:"auth"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;<-:create"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

var user = &UserModel{}

func (AuthModel) TableName() string {
	return "user_auth"
}

func (AuthModel) SetUser(userInfo *UserModel) {
	user = userInfo
}

func User() *UserModel {
	return user
}

func Active() bool {
	if user.Id != "" {
		return true
	}

	return false
}
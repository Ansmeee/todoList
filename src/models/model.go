package models

import (
	"gorm.io/gorm"
)

type Model struct {
	Id        string         `json:"id" gorm:"column:uid" uri:"id" form:"id"`
	CreatedAt string         `json:"created_at" gorm:"autoCreateTime" `
	UpdatedAt string         `json:"updated_at" gorm:"autoUpdateTime" `
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

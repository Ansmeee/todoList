package models

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	Id        int            `json:"id" gorm:"column:uid;<-:create" uri:"id" form:"id"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime;<-:create" `
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime" `
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedAt time.Time      `gorm:"column:created_at; autoCreateTime; <-:create" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at; index" json:"-"`
}

type Task struct {
	ID       int    `gorm:"column:id; primary_key; not null" json:"id"`
	Name     string `gorm:"column:name; type:varchar(255); not null; index" json:"name"`
	Schedule string `gorm:"column:schedule; type:varchar(255); not null" json:"schedule"`
	Message  string `gorm:"column:message; type:text" json:"message"`
	Hash     string `gorm:"column:hash; type:varchar(64); uniqueIndex; not null" json:"hash"`
	Active   bool   `gorm:"column:active; not null; default:true" json:"active"`
	BaseModel
}

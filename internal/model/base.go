package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedAt time.Time      `gorm:"column:created_at; autoCreateTime; <-:create" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at; autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at; index" json:"-"`
}

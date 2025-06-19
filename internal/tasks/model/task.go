package model

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"type:varchar(255);not null;index" json:"name,omitempty"`
	Schedule string `gorm:"type:varchar(255);not null" json:"schedule,omitempty"`
	Message  string `gorm:"type:text" json:"message,omitempty"`
	Hash     string `gorm:"type:varchar(64);uniqueIndex;not null" json:"hash,omitempty"`
	Active   bool   `gorm:"default:true;not null" json:"active,omitempty"`

	CreatedAt time.Time      `gorm:"autoCreateTime;<-:create" json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

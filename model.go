package main

import (
	"gorm.io/gorm"
)

type TaskModel struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null;index" json:"name"` // Có index để tìm kiếm nhanh
	Schedule  string         `gorm:"type:varchar(255);not null" json:"schedule"`
	Message   string         `gorm:"type:text" json:"message"`
	Hash      string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"hash"` // unique và không null
	Active    bool           `gorm:"default:true;not null" json:"active"`               // default true
	CreatedAt int64          `gorm:"autoCreateTime:milli" json:"created_at"`            // timestamp millisecond
	UpdatedAt int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // soft delete, không expose JSON
}

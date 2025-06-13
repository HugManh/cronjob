package main

import (
	"gorm.io/gorm"
)

type TaskModel struct {
	gorm.Model
	Name     string `gorm:"uniqueIndex"`
	Schedule string
	Message  string
	Hash     string
	Active   bool // dùng để tạm dừng chạy task nếu muốn
}

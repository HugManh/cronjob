package main

import (
	"gorm.io/gorm"
)

type TaskModel struct {
	gorm.Model
	Name     string
	Schedule string
	Message  string
	Hash     string `gorm:"uniqueIndex"`
	Active   bool   // dùng để tạm dừng chạy task nếu muốn
}

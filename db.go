package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect DB: %v", err)
	}

	// Auto migrate bảng TaskModel
	if err := db.AutoMigrate(&TaskModel{}); err != nil {
		log.Fatalf("migrate failed: %v", err)
	}

	return db
}

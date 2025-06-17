package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	db = db.Debug()
	if err != nil {
		log.Fatalf("failed to connect DB: %v", err)
	}

	// Auto migrate bảng TaskModel
	// if err := db.AutoMigrate(&TaskModel{}); err != nil {
	// 	log.Fatalf("migrate failed: %v", err)
	// }

	return db
}

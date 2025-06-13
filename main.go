package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config := LoadConfig()
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName, "require",
	)
	db := InitDB(dsn)

	taskManager := NewTaskManager(db)
	// Load task từ DB
	if err := taskManager.LoadTasksFromDB(); err != nil {
		log.Fatalf("Load tasks failed: %v", err)
	}
	taskManager.Cron.Start()

	r := gin.Default()
	RegisterRoutes(r, taskManager)

	r.Run(":8080") // listen on localhost:8080
}

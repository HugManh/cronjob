package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/HugManh/cronjob/internal/routing"
	"github.com/HugManh/cronjob/pkg/config"
	"github.com/HugManh/cronjob/pkg/db"
	"github.com/HugManh/cronjob/pkg/taskmanager"
)

func main() {
	config := config.LoadConfig()
	// Init DB
	db, err := db.Init(config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName, config.DBSsl)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Init TaskManager
	tm := taskmanager.NewTaskManager()
	tm.Cron.Start()

	// Register tasks from DB
	if err := tm.LoadTasksFromDB(db); err != nil {
		log.Fatalf("Registering tasks from DB failed: %v", err)
	}

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		log.Println("[ACCESS] Request >>", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})
	if err := routing.RegisterRoutes(routing.ServerData{DB: db, TaskManager: tm, Router: router}); err != nil {
		fmt.Printf("registering routes: %v\n", err)
		os.Exit(1)
	}

	router.Run(":8080")
}

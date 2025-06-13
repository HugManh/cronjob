package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	taskManager := NewTaskManager()
	taskManager.Cron.Start()

	r := gin.Default()
	RegisterRoutes(r, taskManager)

	r.Run(":8080") // listen on localhost:8080
}

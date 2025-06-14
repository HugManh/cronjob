package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddTaskRequest struct {
	Name     string `json:"name"`
	Schedule string `json:"schedule"` // vd: "*/10 * * * * *"
	Message  string `json:"message"`
}

func RegisterRoutes(r *gin.Engine, tm *TaskManager) {
	v1 := r.Group("/api/v1")
	tasks := v1.Group("/tasks")
	tasks.POST("/", func(c *gin.Context) {
		var req AddTaskRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		id, err := tm.AddTask(req.Name, req.Schedule, req.Message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "task added", "id": id})
	})

	tasks.GET("/", func(c *gin.Context) {
		tasks := tm.GetTasks()
		c.JSON(http.StatusOK, tasks)
	})
}

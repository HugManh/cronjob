package handler

import (
	"fmt"
	"net/http"

	dto_tasks "github.com/HugManh/cronjob/internal/dto/tasks"
	"github.com/HugManh/cronjob/internal/service"
	"github.com/HugManh/cronjob/pkg/https"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	svcTask  *service.TaskService
	svcSlack *service.SlackService
}

func NewTaskHandler(svcTask *service.TaskService, svcSlack *service.SlackService) *TaskHandler {
	return &TaskHandler{svcTask: svcTask, svcSlack: svcSlack}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req dto_tasks.AddTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request",
			"error":   err,
		})
		return
	}

	id, err := h.svcTask.AddTask(req.Name, req.Execute, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task added", "id": id})
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	params := https.ParseQueryParams(c)
	tasks, total, err := h.svcTask.GetTasks(params)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tasks not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":       tasks,
		"page":       params.Page,
		"limit":      params.Limit,
		"total":      total,
		"totalPages": (total + int64(params.Limit) - 1) / int64(params.Limit),
	})
}

func (h *TaskHandler) GetTaskById(c *gin.Context) {
	id := c.Param("id")
	task, err := h.svcTask.GetTaskById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req dto_tasks.AddTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.svcTask.UpdateTask(id, req.Name, req.Execute, req.Message, bool(*req.Active)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task updated", "id": id})
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svcTask.DeleteTask(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task deleted", "id": id})
}

func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		Active *bool `json:"active" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.svcTask.SetTaskActiveStatus(id, *body.Active); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	status := "disabled"
	if *body.Active {
		status = "enabled"
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("task %s successfully", status),
		"id":      id,
	})
}

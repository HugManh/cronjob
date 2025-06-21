package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/HugManh/cronjob/internal/common"
	"github.com/HugManh/cronjob/internal/tasks/dto"
	"github.com/HugManh/cronjob/internal/tasks/service"
)

type TaskHandler struct {
	svc *service.TaskService
}

func NewTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: s}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req dto.AddTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request",
			"error":   err,
		})
		return
	}

	id, err := h.svc.AddTask(req.Name, req.Schedule, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task added", "id": id})
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	params := common.ParseQueryParams(c)
	tasks, total, err := h.svc.GetTasks(params)
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
	task, err := h.svc.GetTaskById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req dto.AddTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.svc.UpdateTask(id, req.Name, req.Schedule, req.Message, bool(*req.Active)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task updated", "id": id})
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteTask(id); err != nil {
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

	if err := h.svc.SetTaskActiveStatus(id, *body.Active); err != nil {
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

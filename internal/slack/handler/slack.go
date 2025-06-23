package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/HugManh/cronjob/internal/common/request"
	"github.com/HugManh/cronjob/internal/slack/dto"
	"github.com/HugManh/cronjob/internal/slack/service"
)

// SlackHandler handles Slack-related HTTP requests
type SlackHandler struct {
	service *service.SlackService
}

// NewSlackHandler creates a new SlackHandler
func NewSlackHandler(s *service.SlackService) *SlackHandler {
	return &SlackHandler{service: s}
}

// Create handles the creation of a new Slack configuration
func (h *SlackHandler) Create(c *gin.Context) {
	var req dto.AddSlackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.CreateSlack(req.BotToken, req.ChatID); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{"message": "Slack configuration created successfully"})
}

// GetSlacks handles retrieving all Slack configurations
func (h *SlackHandler) GetSlacks(c *gin.Context) {
	params := request.ParseQueryParams(c)
	slacks, total, err := h.service.GetSlacks(params)
	if err != nil {
		c.JSON(404, gin.H{"error": "Slacks not found"})
		return
	}

	c.JSON(200, gin.H{
		"data":  slacks,
		"total": total,
	})
}

// GetSlackByID handles retrieving a Slack configuration by ID
func (h *SlackHandler) GetSlackByID(c *gin.Context) {
	id := c.Param("id")
	slack, err := h.service.GetSlackById(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Slack not found"})
		return
	}

	c.JSON(200, slack)
}

// Delete handles deleting a Slack configuration by ID
// func (h *SlackHandler) Delete(c *gin.Context) {
//     id := c.Param("id")
//     if err := h.service.DeleteSlack(id); err != nil {
//         server.HandleError(c, err)
//         return
//     }

//     c.JSON(204, nil) // No content
// }

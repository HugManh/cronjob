package slack

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/HugManh/cronjob/internal/slack/handler"
	"github.com/HugManh/cronjob/internal/slack/repository"
	"github.com/HugManh/cronjob/internal/slack/service"
)

// RegisterRoutes registers the Slack routes with the given Gin router
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB) {
    repo := repository.NewSlackRepo(db)
    slackService := service.NewService(repo)
	h := handler.NewSlackHandler(slackService)

	slackGroup := router.Group("/slacks")
	{
		slackGroup.POST("/", h.Create)
		slackGroup.GET("/", h.GetSlacks)
		slackGroup.GET("/:id", h.GetSlackByID)
	}
}

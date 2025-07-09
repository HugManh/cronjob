package tasks

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	slackRepo "github.com/HugManh/cronjob/internal/slack/repository"
	slackService "github.com/HugManh/cronjob/internal/slack/service"
	"github.com/HugManh/cronjob/internal/tasks/handler"
	taskRepo "github.com/HugManh/cronjob/internal/tasks/repository"
	taskService "github.com/HugManh/cronjob/internal/tasks/service"
	"github.com/HugManh/cronjob/pkg/taskmanager"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, tm *taskmanager.TaskManager) {
	repoTask := taskRepo.NewTaskRepo(db)
	repoSlack := slackRepo.NewSlackRepo(db)
	svcTask := taskService.NewService(repoTask, tm)
	svcSlack := slackService.NewService(repoSlack)
	h := handler.NewTaskHandler(svcTask, svcSlack)

	group := rg.Group("/tasks")
	group.POST("/", h.Create)
	group.GET("/", h.GetTasks)
	group.GET("/:id", h.GetTaskById)
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
	group.POST("/:id/active", h.UpdateStatus)
}

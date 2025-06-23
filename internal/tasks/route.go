package tasks

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/HugManh/cronjob/internal/tasks/handler"
	"github.com/HugManh/cronjob/internal/tasks/repository"
	"github.com/HugManh/cronjob/internal/tasks/service"
	"github.com/HugManh/cronjob/pkg/taskmanager"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, tm *taskmanager.TaskManager) {
	repo := repository.NewTaskRepo(db)
	svc := service.NewService(repo, tm)
	h := handler.NewTaskHandler(svc)

	group := rg.Group("/tasks")
	group.POST("/", h.Create)
	group.GET("/", h.GetTasks)
	group.GET("/:id", h.GetTaskById)
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
	group.POST("/:id/active", h.UpdateStatus)
}

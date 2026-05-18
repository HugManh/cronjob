package routing

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/HugManh/cronjob/internal/handler"
	"github.com/HugManh/cronjob/internal/repository"
	"github.com/HugManh/cronjob/internal/service"
)

func registerAPIRoutes(rg *gin.RouterGroup, db *gorm.DB, tm *service.TaskManager) {
	registerTaskAPIRoutes(rg, db, tm)
	registerSlackAPIRoutes(rg, db)
}

func registerTaskAPIRoutes(rg *gin.RouterGroup, db *gorm.DB, tm *service.TaskManager) {
	taskRepo := repository.NewTaskRepository(db)
	slackRepo := repository.NewSlackRepository(db)
	taskService := service.NewTaskService(taskRepo, tm)
	slackService := service.NewSlackService(slackRepo)
	taskHandler := handler.NewTaskHandler(taskService, slackService)

	tasks := rg.Group("/tasks")
	tasks.POST("/", taskHandler.Create)
	tasks.GET("/", taskHandler.GetTasks)
	tasks.GET("/:id", taskHandler.GetTaskByID)
	tasks.GET("/:id/logs", taskHandler.GetTaskLogs)
	tasks.PUT("/:id", taskHandler.Update)
	tasks.DELETE("/:id", taskHandler.Delete)
	tasks.POST("/:id/active", taskHandler.UpdateStatus)
}

func registerSlackAPIRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	slackRepo := repository.NewSlackRepository(db)
	slackService := service.NewSlackService(slackRepo)
	slackHandler := handler.NewSlackHandler(slackService)

	slacks := rg.Group("/slacks")
	slacks.POST("/", slackHandler.Create)
	slacks.GET("/", slackHandler.GetSlacks)
	slacks.GET("/:id", slackHandler.GetSlackByID)
}

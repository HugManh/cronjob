package routing

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/HugManh/cronjob/internal/handler"
	"github.com/HugManh/cronjob/internal/repository"
	"github.com/HugManh/cronjob/internal/service"
	view "github.com/HugManh/cronjob/internal/web"
)

type ServerData struct {
	DB          *gorm.DB
	Router      *gin.Engine
	TaskManager *service.TaskManager
}

func registerRoutes(ser *ServerData) error {

	r := ser.Router
	vg := r.Group("view")
	view.RegisterRoutes(vg, ser.DB, ser.TaskManager)

	// API routes
	api := r.Group("/api/v1")
	Routes1(api, ser.DB, ser.TaskManager)
	Routes2(api, ser.DB)
	return nil
}

func Routes1(rg *gin.RouterGroup, db *gorm.DB, tm *service.TaskManager) {
	repoTask := repository.NewTaskRepo(db)
	repoSlack := repository.NewSlackRepo(db)
	svcTask := service.NewService2(repoTask, tm)
	svcSlack := service.NewService1(repoSlack)
	h := handler.NewTaskHandler(svcTask, svcSlack)

	group := rg.Group("/tasks")
	group.POST("/", h.Create)
	group.GET("/", h.GetTasks)
	group.GET("/:id", h.GetTaskById)
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
	group.POST("/:id/active", h.UpdateStatus)
}

// RegisterRoutes registers the Slack routes with the given Gin router
func Routes2(router *gin.RouterGroup, db *gorm.DB) {
	repo := repository.NewSlackRepo(db)
	slackService := service.NewService1(repo)
	h := handler.NewSlackHandler(slackService)

	slackGroup := router.Group("/slacks")
	{
		slackGroup.POST("/", h.Create)
		slackGroup.GET("/", h.GetSlacks)
		slackGroup.GET("/:id", h.GetSlackByID)
	}
}

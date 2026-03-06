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

	// User UI routes (Served from root)
	r.StaticFile("/", "./public/index.html")
	r.StaticFile("/style.css", "./public/style.css")
	r.StaticFile("/app.js", "./public/app.js")

	// View routes
	vg := r.Group("view")
	view.RegisterRoutes(vg, ser.DB, ser.TaskManager)

	// API routes
	api := r.Group("/api/v1")
	Routes1(api, ser.DB, ser.TaskManager)
	Routes2(api, ser.DB)

	// 404 handler for routes that don't exist
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error": "Route not found",
			"path":  c.Request.URL.Path,
		})
	})

	// 405 handler for method not allowed
	r.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{
			"error":  "Method not allowed",
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		})
	})

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
	group.GET("/:id/logs", h.GetTaskLogs)
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

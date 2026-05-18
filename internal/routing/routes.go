package routing

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/HugManh/cronjob/internal/service"
)

type ServerData struct {
	DB          *gorm.DB
	Router      *gin.Engine
	TaskManager *service.TaskManager
}

func registerRoutes(server *ServerData) error {
	registerWebRoutes(server.Router, server.DB, server.TaskManager)
	registerAPIRoutes(server.Router.Group("/api/v1"), server.DB, server.TaskManager)
	registerErrorHandlers(server.Router)
	return nil
}

func registerErrorHandlers(router *gin.Engine) {
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error": "Route not found",
			"path":  c.Request.URL.Path,
		})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{
			"error":  "Method not allowed",
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		})
	})
}

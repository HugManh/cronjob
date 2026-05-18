package routing

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/HugManh/cronjob/internal/service"
	view "github.com/HugManh/cronjob/internal/web"
)

func registerWebRoutes(router *gin.Engine, db *gorm.DB, tm *service.TaskManager) {
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/view/tasks")
	})

	view.RegisterRoutes(router.Group("/view"), db, tm)
}

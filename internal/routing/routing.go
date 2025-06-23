package routing

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/HugManh/cronjob/internal/slack"
	"github.com/HugManh/cronjob/internal/tasks"
	view "github.com/HugManh/cronjob/internal/web"
	"github.com/HugManh/cronjob/pkg/taskmanager"
)

type ServerData struct {
	DB          *gorm.DB
	Router      *gin.Engine
	TaskManager *taskmanager.TaskManager
}

func RegisterRoutes(ser ServerData) error {

	r := ser.Router
	vg := r.Group("view")
	view.RegisterRoutes(vg, ser.DB, ser.TaskManager)

	// API routes
	api := r.Group("/api/v1")
	tasks.RegisterRoutes(api, ser.DB, ser.TaskManager)
	slack.RegisterRoutes(api, ser.DB)
	return nil
}

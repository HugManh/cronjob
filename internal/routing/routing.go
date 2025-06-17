package routing

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/HugManh/cronjob/internal/tasks"
	"github.com/HugManh/cronjob/pkg/taskmanager"
)

type ServerData struct {
	DB          *gorm.DB
	Router      *gin.Engine
	TaskManager *taskmanager.TaskManager
}

func RegisterRoutes(ser ServerData) error {
	api := ser.Router.Group("/api/v1")
	tasks.RegisterRoutes(api, ser.DB, ser.TaskManager)
	return nil
}

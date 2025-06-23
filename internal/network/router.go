package network

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/HugManh/cronjob/internal/routing"
	"github.com/HugManh/cronjob/pkg/taskmanager"
	"github.com/gin-gonic/gin"
)

type BaseRouter interface {
	GetEngine() *gin.Engine
	Start(ip string, port uint16)
}

type Router interface {
	BaseRouter
	LoadControllers(*gorm.DB, *taskmanager.TaskManager)
}

type router struct {
	engine *gin.Engine
}

func NewRouter(mode string) Router {
	gin.SetMode(mode)
	eng := gin.Default()
	r := router{
		engine: eng,
	}
	return &r
}

func (r *router) LoadControllers(db *gorm.DB, tm *taskmanager.TaskManager) {
	r.engine.Use(func(c *gin.Context) {
		log.Println("[ACCESS] Request >>", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})
	if err := routing.RegisterRoutes(
		routing.ServerData{
			DB:          db,
			TaskManager: tm,
			Router:      r.engine,
		}); err != nil {
		log.Errorf("registering routes: %v\n", err)
		os.Exit(1)
	}
}

func (r *router) GetEngine() *gin.Engine {
	return r.engine
}

func (r *router) Start(ip string, port uint16) {
	address := fmt.Sprintf("%s:%d", ip, port)
	r.engine.Run(address)
}

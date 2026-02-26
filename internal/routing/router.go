package routing

import (
	"fmt"
	"os"

	"github.com/HugManh/cronjob/internal/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BaseRouter interface {
	GetEngine() *gin.Engine
	Start(ip string, port uint16)
}

type Router interface {
	BaseRouter
	LoadControllers(*gorm.DB, *service.TaskManager)
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

func (r *router) LoadControllers(db *gorm.DB, tm *service.TaskManager) {
	r.engine.Use(func(c *gin.Context) {
		log.Println("[ACCESS] Request >>", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})
	if err := registerRoutes(&ServerData{
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

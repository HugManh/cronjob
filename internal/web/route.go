package view

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/CloudyKit/jet/v6"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/HugManh/cronjob/internal/repository"
	"github.com/HugManh/cronjob/internal/service"
	"github.com/HugManh/cronjob/pkg/httpx"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./views"),
	jet.DevelopmentMode(true), // remove or set false in production
)

func Init() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Working directory:", cwd)
}

// Render HTML views
func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, tm *service.TaskManager) {
	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo, tm)
	slackRepo := repository.NewSlackRepository(db)
	slackService := service.NewSlackService(slackRepo)

	tasks := rg.Group("/tasks")
	tasks.GET("/", func(c *gin.Context) {
		tmpl, err := views.GetTemplate("tasks/list")
		if err != nil {
			c.String(http.StatusInternalServerError, "template error: %v", err)
			return
		}

		err = tmpl.Execute(c.Writer, nil, nil)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("render error: %v", err))
		}
	})
	tasks.GET("/items", func(c *gin.Context) {
		params := httpx.ParseQueryParams(c)
		tasks, _, err := taskService.GetTasks(params)
		if err != nil {
			c.String(http.StatusNotFound, "tasks not found: %v", err)
			return
		}

		tmpl, err := views.GetTemplate("tasks/items")
		if err != nil {
			c.String(http.StatusInternalServerError, "template error: %v", err)
			return
		}

		vars := make(jet.VarMap)
		vars.Set("tasks", tasks)

		err = tmpl.Execute(c.Writer, vars, tasks)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("render error: %v", err))
		}
	})
	tasks.GET("/new", func(c *gin.Context) {
		tmpl, err := views.GetTemplate("tasks/new")
		if err != nil {
			c.String(http.StatusInternalServerError, "template error: %v", err)
			return
		}

		tmpl.Execute(c.Writer, nil, nil)
	})
	tasks.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		task, err := taskService.GetTaskByID(id)
		if err != nil {
			c.String(http.StatusNotFound, "task not found: %v", err)
			return
		}

		tmpl, err := views.GetTemplate("tasks/detail")
		if err != nil {
			c.String(http.StatusInternalServerError, "template error: %v", err)
			return
		}

		vars := make(jet.VarMap)
		vars.Set("task", task)

		err = tmpl.Execute(c.Writer, vars, task)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("render error: %v", err))
		}
	})

	slacks := rg.Group("/slacks")
	slacks.GET("/", func(c *gin.Context) {
		tmpl, err := views.GetTemplate("slacks/list")
		if err != nil {
			c.String(http.StatusInternalServerError, "template error: %v", err)
			return
		}

		if err := tmpl.Execute(c.Writer, nil, nil); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("render error: %v", err))
		}
	})
	slacks.GET("/items", func(c *gin.Context) {
		params := httpx.ParseQueryParams(c)
		slacks, _, err := slackService.GetSlacks(params)
		if err != nil {
			c.String(http.StatusNotFound, "slacks not found: %v", err)
			return
		}

		tmpl, err := views.GetTemplate("slacks/items")
		if err != nil {
			c.String(http.StatusInternalServerError, "template error: %v", err)
			return
		}

		vars := make(jet.VarMap)
		vars.Set("slacks", slacks)

		if err := tmpl.Execute(c.Writer, vars, slacks); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("render error: %v", err))
		}
	})

}

package view

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/CloudyKit/jet/v6"
	"github.com/HugManh/cronjob/internal/common/request"
	repository1 "github.com/HugManh/cronjob/internal/slack/repository"
	service1 "github.com/HugManh/cronjob/internal/slack/service"
	repository "github.com/HugManh/cronjob/internal/tasks/repository"
	service "github.com/HugManh/cronjob/internal/tasks/service"
	"github.com/HugManh/cronjob/pkg/taskmanager"
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
func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, tm *taskmanager.TaskManager) {
	repo := repository.NewTaskRepo(db)
	svc := service.NewService(repo, tm)
	repo1 := repository1.NewSlackRepo(db)
	svc1 := service1.NewService(repo1)

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
		params := request.ParseQueryParams(c)
		tasks, _, err := svc.GetTasks(params)
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
		task, err := svc.GetTaskById(id)
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

		err = tmpl.Execute(c.Writer, nil, nil)
		fmt.Printf("render error: %v", err)
	})
	slacks.GET("/items", func(c *gin.Context) {
		params := request.ParseQueryParams(c)
		slacks, _, err := svc1.GetSlacks(params)
		if err != nil {
			c.String(http.StatusNotFound, "slacks not found: %v", err)
			return
		}

		fmt.Printf("---------- %+v", slacks)

		tmpl, err := views.GetTemplate("slacks/items")
		if err != nil {
			c.String(http.StatusInternalServerError, "template error: %v", err)
			return
		}

		vars := make(jet.VarMap)
		vars.Set("slacks", slacks)

		err = tmpl.Execute(c.Writer, vars, slacks)
        fmt.Printf("render error: %v", err)
	})

}

package view

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/CloudyKit/jet/v6"
	"github.com/HugManh/cronjob/internal/tasks/repository"
	"github.com/HugManh/cronjob/internal/tasks/service"
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
	fmt.Println("Working directory:", cwd)
}

// Render HTML views
func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, tm *taskmanager.TaskManager) {
	repo := repository.NewTaskRepo(db)
	svc := service.NewTaskService(repo, tm)

	group := rg.Group("/tasks")
	group.GET("/", func(c *gin.Context) {
		fmt.Println(">>> List")
		tasks, _ := svc.GetTasks()

		tmpl, err := views.GetTemplate("tasks/list.jet")
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
	group.GET("/new", func(c *gin.Context) {
		fmt.Println(">>> New")
		tmpl, err := views.GetTemplate("tasks/new.jet")
		if err != nil {
			c.String(http.StatusInternalServerError, "template error: %v", err)
			return
		}

		tmpl.Execute(c.Writer, nil, nil)
	})
	group.GET("/:id", func(c *gin.Context) {
		fmt.Println(">>> id:")
		id := c.Param("id")
		task, err := svc.GetTaskById(id)
		if err != nil {
			c.String(http.StatusNotFound, "task not found: %v", err)
			return
		}

		tmpl, err := views.GetTemplate("tasks/detail.jet")
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

}

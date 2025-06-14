package main

import (
	"fmt"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/gin-gonic/gin"
)

type AddTaskRequest struct {
	Name     string `json:"name"`
	Schedule string `json:"schedule"` // vd: "*/10 * * * * *"
	Message  string `json:"message"`
}

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./views"),
	jet.DevelopmentMode(true), // remove or set false in production
)

func RegisterRoutes(r *gin.Engine, tm *TaskManager) {
	v1 := r.Group("/api/v1")
	tasks := v1.Group("/tasks")
	tasks.POST("/", func(c *gin.Context) {
		var req AddTaskRequest
		// if err := c.ShouldBindJSON(&req); err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		// 	return
		// }

		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		id, err := tm.AddTask(req.Name, req.Schedule, req.Message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "task added", "id": id})
	})

	tasks.GET("/", func(c *gin.Context) {
		tasks := tm.GetTasks()
		c.JSON(http.StatusOK, tasks)
	})

	view := r.Group("/view")
	viewTasks := view.Group("/tasks")
	// viewTasks.GET("/new", func(c *gin.Context) {
	// 	tmpl, err := views.GetTemplate("tasks/new.jet")
	// 	if err != nil {
	// 		c.String(http.StatusInternalServerError, "template error: %v", err)
	// 		return
	// 	}

	// 	tmpl.Execute(c.Writer, nil, nil)
	// })

	viewTasks.GET("/", func(c *gin.Context) {
		tasks := tm.GetTasks()

		tmpl, err := views.GetTemplate("tasks/list.jet")
		if err != nil {
			c.String(http.StatusInternalServerError, "template error: %v", err)
			return
		}

		// vars := make(jet.VarMap)
		// vars.Set("tasks", tasks)

		fmt.Println("Tasks list:", tasks)

		err = tmpl.Execute(c.Writer, nil, tasks)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("render error: %v", err))
		}
	})
}

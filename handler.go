package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/gin-gonic/gin"
)

type BoolString bool

func (b *BoolString) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), "\"")
	switch str {
	case "true", "True", "1":
		*b = true
	case "false", "False", "0":
		*b = false
	default:
		return fmt.Errorf("invalid boolean: %s", str)
	}
	return nil
}

type AddTaskRequest struct {
	Name     string     `json:"name"`
	Schedule string     `json:"schedule"` // vd: "*/10 * * * * *"
	Message  string     `json:"message"`
	Active   BoolString `json:"active"` // Trạng thái kích hoạt
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
		if err := c.ShouldBindJSON(&req); err != nil {
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

	tasks.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		task, err := tm.GetTaskById(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusOK, task)
	})

	tasks.PUT("/:id", func(c *gin.Context) {
		id := c.Param("id")
		var req AddTaskRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		if err := tm.UpdateTask(id, req.Name, req.Schedule, req.Message, bool(req.Active)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "task updated", "id": id})
	})

	tasks.POST("/:id/active", func(c *gin.Context) {
		id := c.Param("id")

		var body struct {
			Active *bool `json:"active" binding:"required"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		if err := tm.SetTaskActiveStatus(id, *body.Active); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		status := "disabled"
		if *body.Active {
			status = "enabled"
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("task %s successfully", status),
			"id":      id,
		})
	})

	// Render HTML views
	view := r.Group("/view")
	viewTasks := view.Group("/tasks")
	viewTasks.GET("/new", func(c *gin.Context) {
		tmpl, err := views.GetTemplate("tasks/new.jet")
		if err != nil {
			c.String(http.StatusInternalServerError, "template error: %v", err)
			return
		}

		tmpl.Execute(c.Writer, nil, nil)
	})
	viewTasks.GET("/", func(c *gin.Context) {
		tasks := tm.GetTasks()

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

	viewTasks.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		task, err := tm.GetTaskById(id)
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

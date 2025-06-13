package main

import (
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
)

type Task struct {
	ID       cron.EntryID
	Name     string
	Schedule string
	Message  string
}

type TaskManager struct {
	Cron  *cron.Cron
	Tasks map[cron.EntryID]Task
	mu    sync.Mutex
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		Cron:  cron.New(cron.WithSeconds()),
		Tasks: make(map[cron.EntryID]Task),
	}
}

func (tm *TaskManager) AddTask(name, schedule, message string) (cron.EntryID, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, task := range tm.Tasks {
		if task.Name == name {
			return 0, fmt.Errorf("task name '%s' already exists", name)
		}
	}

	id, err := tm.Cron.AddFunc(schedule, func() {
		fmt.Printf("[TASK %s] %s\n", name, message)
	})
	if err != nil {
		return 0, err
	}

	tm.Tasks[id] = Task{
		ID:       id,
        Name:     name,
		Schedule: schedule,
		Message:  message,
	}

	return id, nil
}

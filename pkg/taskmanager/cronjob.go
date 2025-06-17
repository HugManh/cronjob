package taskmanager

import (
	"fmt"
	"log"
	"sync"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Task struct {
	ID       cron.EntryID
	Hash     string
	Name     string
	Schedule string
	Message  string
}

type TaskManager struct {
	DB    *gorm.DB
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

// Register a new task in the cron scheduler
func (tm *TaskManager) RegisterTask(hash, name, schedule, message string) (cron.EntryID, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Check trùng tên
	for _, task := range tm.Tasks {
		if task.Hash == hash {
			return 0, fmt.Errorf("task name '%s' already exists", name)
		}
	}

	id, err := tm.Cron.AddFunc(schedule, func() {
		log.Printf("[TASK %s][%s] %s", hash, name, message)
	})
	if err != nil {
		return 0, err
	}

	tm.Tasks[id] = Task{
		ID:       id,
		Hash:     hash,
		Name:     name,
		Schedule: schedule,
		Message:  message,
	}

	log.Printf("✅ Registered task: %s | %s", name, schedule)
	return id, nil
}

// Remove task in the cron scheduler
func (tm *TaskManager) RemoveTaskFromCronByHash(taskHash string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for id, t := range tm.Tasks {
		if t.Hash == taskHash { //
			log.Printf("Checking task: %v against %s\n", t, taskHash)
			tm.Cron.Remove(id)
			delete(tm.Tasks, id)
			log.Printf("✅ Removed task with hash %s from cron", taskHash)
			return nil
		}
	}

	log.Printf("⚠️ Task with hash %s not found in cron", taskHash)
	return fmt.Errorf("task with hash %s not found in cron", taskHash)
}

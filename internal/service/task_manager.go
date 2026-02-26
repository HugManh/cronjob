package service

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/HugManh/cronjob/internal/model"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Task struct {
	ID      cron.EntryID
	Hash    string
	Name    string
	Execute string
	Message string
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

const invalid_task = "000001"

func (tm *TaskManager) Startup(db *gorm.DB) error {
	var tasks []model.Task
	if err := db.Where("active = ?", true).Find(&tasks).Error; err != nil {
		log.Errorf("Failed to load tasks from DB: %v", err)
		return err
	}

	for _, t := range tasks {
		log.Infof("Loading task from DB: Name: %s Execute: %s Message: %s Hash: %s", t.Name, t.Execute, t.Message, t.Hash)
		_, err := tm.RegisterTask(t.Hash, t.Name, t.Execute, t.Message)
		if err != nil {
			log.Infof("Failed to add task %s: %v", t.Name, err)
			// Nếu task không hợp lệ, đánh dấu nó là inactive để tránh lỗi khi khởi động lại
			if err := db.Model(&model.Task{}).
				Where("id = ?", t.ID).
				Updates(map[string]interface{}{
					"active": false,
					"code":   invalid_task,
				}).Error; err != nil {
				log.Errorf("Failed to update task %s to inactive: %v", t.Name, err)
			}
		}
	}

	return nil
}

// Register a new task in the cron scheduler
func (tm *TaskManager) RegisterTask(hash, name, execute, message string) (cron.EntryID, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Check trùng tên
	for _, task := range tm.Tasks {
		if task.Hash == hash {
			return 0, fmt.Errorf("task name '%s' already exists", name)
		}
	}

	id, err := tm.Cron.AddFunc(execute, func() {
		log.Infof("[TASK %s][%s] %s", hash, name, message)
	})
	if err != nil {
		return 0, err
	}

	tm.Tasks[id] = Task{
		ID:      id,
		Hash:    hash,
		Name:    name,
		Execute: execute,
		Message: message,
	}

	log.Infof("✅ Registered task: %s | %s", name, execute)
	return id, nil
}

// Remove task in the cron scheduler
func (tm *TaskManager) RemoveTaskFromCronByHash(taskHash string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for id, t := range tm.Tasks {
		if t.Hash == taskHash {
			log.Infof("Checking task: %v against %s\n", t, taskHash)
			tm.Cron.Remove(id)
			delete(tm.Tasks, id)
			log.Infof("✅ Removed task with hash %s from cron", taskHash)
			return nil
		}
	}

	log.Infof("⚠️ Task with hash %s not found in cron", taskHash)
	return fmt.Errorf("task with hash %s not found in cron", taskHash)
}

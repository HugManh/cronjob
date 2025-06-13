package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
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

func NewTaskManager(db *gorm.DB) *TaskManager {
	return &TaskManager{
		DB:    db,
		Cron:  cron.New(cron.WithSeconds()),
		Tasks: make(map[cron.EntryID]Task),
	}
}

// Load task từ DB và đăng ký lại cron
func (tm *TaskManager) LoadTasksFromDB() error {
	var taskModels []TaskModel
	if err := tm.DB.Where("active = ?", true).Find(&taskModels).Error; err != nil {
		log.Printf("Failed to load tasks from DB: %v", err)
		return err
	}

	for _, t := range taskModels {
		log.Printf("Loading task from DB: %s, Schedule: %s, Message: %s", t.Name, t.Schedule, t.Message)
		_, err := tm.registerTask(t.Hash, t.Name, t.Schedule, t.Message)
		if err != nil {
			log.Printf("Failed to add task %s: %v", t.Name, err)
		}
	}

	return nil
}

func (tm *TaskManager) registerTask(hash, name, schedule, message string) (cron.EntryID, error) {
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
		Name:     name,
		Schedule: schedule,
		Message:  message,
	}

	log.Printf("✅ Registered task: %s | %s", name, schedule)
	return id, nil
}

func (tm *TaskManager) AddTask(name, schedule, message string) (cron.EntryID, error) {
	// Đăng ký task
	hash := uuid.New().String()
	id, err := tm.registerTask(hash, name, schedule, message)
	if err != nil {
		return 0, err
	}

	// Lưu vào DB (active = true)
	taskModel := TaskModel{
		Name:     name,
		Schedule: schedule,
		Message:  message,
		Hash:     hash,
		Active:   true,
	}
	if err := tm.DB.Create(&taskModel).Error; err != nil {
		// Nếu lưu DB lỗi thì rollback cron luôn
		tm.Cron.Remove(id)
		delete(tm.Tasks, id)
		return 0, err
	}

	return id, nil
}

func (tm *TaskManager) DisableTaskByName(hash string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	var entryID cron.EntryID = 0
	for id, task := range tm.Tasks {
		if task.Hash == hash {
			entryID = id
			break
		}
	}

	if entryID == 0 {
		return fmt.Errorf("task %s not found", hash)
	}

	tm.Cron.Remove(entryID)
	delete(tm.Tasks, entryID)

	// Cập nhật DB: đặt active = false
	return tm.DB.Model(&TaskModel{}).Where("hash = ?", hash).Update("active", false).Error
}

func (tm *TaskManager) DeleteTaskByName(hash string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	var entryID cron.EntryID
	for id, task := range tm.Tasks {
		if task.Hash == hash {
			entryID = id
			break
		}
	}
	if entryID != 0 {
		tm.Cron.Remove(entryID)
		delete(tm.Tasks, entryID)
	}

	return tm.DB.Where("hash = ?", hash).Delete(&TaskModel{}).Error
}

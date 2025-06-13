package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Task struct {
	ID       cron.EntryID
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
		return err
	}

	for _, t := range taskModels {
		_, err := tm.AddTask(t.Name, t.Schedule, t.Message)
		if err != nil {
			log.Printf("Failed to add task %s: %v", t.Name, err)
		}
	}

	return nil
}

func (tm *TaskManager) AddTask(name, schedule, message string) (cron.EntryID, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Check trùng tên trong map
	for _, task := range tm.Tasks {
		if task.Name == name {
			return 0, fmt.Errorf("task name '%s' already exists", name)
		}
	}

	// Thêm task vào cron
	id, err := tm.Cron.AddFunc(schedule, func() {
		log.Printf("[TASK %s] %s\n", name, message)
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

	// Lưu vào DB (active = true)
	taskModel := TaskModel{
		Name:     name,
		Schedule: schedule,
		Message:  message,
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

func (tm *TaskManager) RemoveTaskByName(name string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	var entryID cron.EntryID = 0
	for id, task := range tm.Tasks {
		if task.Name == name {
			entryID = id
			break
		}
	}

	if entryID == 0 {
		return fmt.Errorf("task %s not found", name)
	}

	tm.Cron.Remove(entryID)
	delete(tm.Tasks, entryID)

	// Cập nhật DB: đặt active = false
	return tm.DB.Model(&TaskModel{}).Where("name = ?", name).Update("active", false).Error
}

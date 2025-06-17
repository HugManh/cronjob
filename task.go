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
	var tasks []TaskModel
	if err := tm.DB.Where("active = ?", true).Find(&tasks).Error; err != nil {
		log.Printf("Failed to load tasks from DB: %v", err)
		return err
	}

	for _, t := range tasks {
		log.Printf("Loading task from DB: Name: %s Schedule: %s Message: %s Hash: %s", t.Name, t.Schedule, t.Message, t.Hash)
		_, err := tm.registerTask(t.Hash, t.Name, t.Schedule, t.Message)
		if err != nil {
			log.Printf("Failed to add task %s: %v", t.Name, err)
		}
	}

	return nil
}

// Register a new task in the cron scheduler
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
		Hash:     hash,
		Name:     name,
		Schedule: schedule,
		Message:  message,
	}

	log.Printf("✅ Registered task: %s | %s", name, schedule)
	return id, nil
}

// AddTask adds a new task to the cron scheduler and saves it to the database
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

func (tm *TaskManager) GetTasks() []TaskModel {
	var tasks []TaskModel
	if err := tm.DB.Find(&tasks).Error; err != nil {
		return []TaskModel{}
	}
	return tasks
}

func (tm *TaskManager) GetTaskById(id string) (*TaskModel, error) {
	var task TaskModel
	if err := tm.DB.Where("id = ?", id).First(&task).Error; err != nil {
		return nil, fmt.Errorf("task with id %s not found: %v", id, err)
	}
	return &task, nil
}

func (tm *TaskManager) SetTaskActiveStatus(id string, active bool) error {
	fmt.Println("----------------------- SetTaskActiveStatus -----------------------", id, active)
	var task TaskModel
	if err := tm.DB.Where("id = ?", id).First(&task).Error; err != nil {
		return fmt.Errorf("task with id %s not found: %v", id, err)
	}

	if task.Active == active {
		fmt.Printf("task with id %s is already in desired state", id)
		return nil
	}

	tm.mu.Lock()
	defer tm.mu.Unlock()

	if active {
		// Thêm lại vào cron
		entryID, err := tm.registerTask(task.Hash, task.Name, task.Schedule, task.Message)
		if err != nil {
			return fmt.Errorf("failed to schedule task: %v", err)
		}
		tm.Tasks[entryID] = Task{Hash: task.Hash}
	} else {
		// Gỡ khỏi cron
		var entryID cron.EntryID
		found := false
		for id, t := range tm.Tasks {
			fmt.Printf("Checking task: %v against %s\n", t, task.Hash)
			if t.Hash == task.Hash {
				entryID = id
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("task with hash %s not found in cron", task.Hash)
		}
		tm.Cron.Remove(entryID)
		delete(tm.Tasks, entryID)
	}
	// Cập nhật DB: đặt active
	return tm.DB.Model(&task).Update("active", active).Error
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

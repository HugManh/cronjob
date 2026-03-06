package service

import (
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/HugManh/cronjob/internal/model"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

const maxLogEntries = 50

// LogEntry represents a single task execution log
type LogEntry struct {
	Time    string `json:"time"`
	Message string `json:"message"`
}

type Task struct {
	ID      cron.EntryID
	Hash    string
	Name    string
	Execute string
	Message string
}

type TaskManager struct {
	DB       *gorm.DB
	Cron     *cron.Cron
	Tasks    map[cron.EntryID]Task
	TaskLogs map[string][]LogEntry // keyed by task hash
	mu       sync.Mutex
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		Cron:     cron.New(cron.WithSeconds()),
		Tasks:    make(map[cron.EntryID]Task),
		TaskLogs: make(map[string][]LogEntry),
	}
}

// AddLog appends a log entry for a task hash, keeping only the last maxLogEntries
func (tm *TaskManager) AddLog(hash, message string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	entry := LogEntry{
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		Message: message,
	}
	logs := tm.TaskLogs[hash]
	logs = append(logs, entry)
	if len(logs) > maxLogEntries {
		logs = logs[len(logs)-maxLogEntries:]
	}
	tm.TaskLogs[hash] = logs
}

// GetLogs returns log entries for a given task hash (most recent first)
func (tm *TaskManager) GetLogs(hash string) []LogEntry {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	logs := tm.TaskLogs[hash]
	// Return a reversed copy so newest is first
	result := make([]LogEntry, len(logs))
	for i, e := range logs {
		result[len(logs)-1-i] = e
	}
	return result
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

// RegisterTask adds a task to the cron scheduler and wires up log capture
func (tm *TaskManager) RegisterTask(hash, name, execute, message string) (cron.EntryID, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Check for duplicate hash
	for _, task := range tm.Tasks {
		if task.Hash == hash {
			return 0, fmt.Errorf("task name '%s' already exists", name)
		}
	}

	id, err := tm.Cron.AddFunc(execute, func() {
		msg := fmt.Sprintf("[TASK %s][%s] %s", hash, name, message)
		log.Info(msg)
		tm.AddLog(hash, message)
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

// RemoveTaskFromCronByHash removes a task from the cron scheduler by its hash
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

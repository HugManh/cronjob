package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"

	"github.com/HugManh/cronjob/internal/common/request"
	"github.com/HugManh/cronjob/internal/tasks/model"
	"github.com/HugManh/cronjob/internal/tasks/repository"
	"github.com/HugManh/cronjob/pkg/taskmanager"
)

type TaskService struct {
	repo *repository.TaskRepo
	tm   *taskmanager.TaskManager
}

func NewTaskService(r *repository.TaskRepo, tm *taskmanager.TaskManager) *TaskService {
	return &TaskService{repo: r, tm: tm}
}

// AddTask adds a new task to the cron scheduler and saves it to the database
func (s *TaskService) AddTask(name, schedule, message string) (cron.EntryID, error) {
	// Đăng ký task
	hash := uuid.New().String()
	id, err := s.tm.RegisterTask(hash, name, schedule, message)
	if err != nil {
		return 0, err
	}

	// Lưu vào DB (active = true)
	task := model.Task{
		Name:     name,
		Schedule: schedule,
		Message:  message,
		Hash:     hash,
		Active:   true,
	}
	if err := s.repo.Create(&task); err != nil {
		// Nếu lưu DB lỗi thì rollback cron luôn
		s.tm.Cron.Remove(id)
		delete(s.tm.Tasks, id)
		return 0, err
	}

	return id, nil
}

func (s *TaskService) GetTasks(params request.QueryParams) ([]model.Task, int64, error) {
	return s.repo.GetAll(params)
}

func (s *TaskService) GetTaskById(id string) (*model.Task, error) {
	return s.repo.GetByID(id)
}

func (s *TaskService) SetTaskActiveStatus(id string, active bool) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("task with id %s not found: %v", id, err)
	}

	if task.Active == active {
		log.Printf("task with id %s is already in desired state", id)
		return nil
	}

	if active {
		// Thêm lại vào cron
		entryID, err := s.tm.RegisterTask(task.Hash, task.Name, task.Schedule, task.Message)
		if err != nil {
			return fmt.Errorf("failed to schedule task: %v", err)
		}
		s.tm.Tasks[entryID] = taskmanager.Task{Hash: task.Hash}
	} else {
		if err := s.tm.RemoveTaskFromCronByHash(task.Hash); err != nil {
			return fmt.Errorf("failed to remove task: %v", err)
		}
	}

	// Cập nhật DB: đặt active
	task.Active = active
	return s.repo.Update(task)
}

func (s *TaskService) UpdateTask(id string, name, schedule, message string, active bool) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("task with id %s not found: %v", id, err)
	}

	// Nếu không có thay đổi, không làm gì
	if task.Name == name && task.Schedule == schedule && task.Message == message && task.Active == active {
		return nil
	}

	if task.Active {
		if err := s.tm.RemoveTaskFromCronByHash(task.Hash); err != nil {
			return fmt.Errorf("failed to remove task: %v", err)
		}
	}

	// Cập nhật DB
	task.Name = name
	task.Schedule = schedule
	task.Message = message
	task.Active = active
	if err := s.repo.Update(task); err != nil {
		return fmt.Errorf("failed to update task in DB: %v", err)
	}

	// Đăng ký lại vào cron nếu đang active
	if task.Active {
		if _, err := s.tm.RegisterTask(task.Hash, name, schedule, message); err != nil {
			return fmt.Errorf("failed to re-register updated task in cron: %v", err)
		}
	}

	log.Printf("✅ Updated task: %s", name)
	return nil
}

func (s *TaskService) DeleteTask(id string) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("task with id %s not found: %v", id, err)
	}

	if task.Active {
		if err := s.tm.RemoveTaskFromCronByHash(task.Hash); err != nil {
			return fmt.Errorf("failed to remove task: %v", err)
		}
	}

	return s.repo.Delete(id)
}

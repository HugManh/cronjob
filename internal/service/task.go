package service

import (
	"fmt"

	"github.com/HugManh/cronjob/internal/model"
	"github.com/HugManh/cronjob/internal/repository"
	"github.com/HugManh/cronjob/pkg/httpx"
	"github.com/HugManh/cronjob/pkg/logger"
	"github.com/google/uuid"
)

type TaskService struct {
	repo *repository.TaskRepository
	tm   *TaskManager
}

func NewTaskService(r *repository.TaskRepository, tm *TaskManager) *TaskService {
	return &TaskService{repo: r, tm: tm}
}

// AddTask creates a task, registers it with cron, and returns the persisted task.
func (s *TaskService) AddTask(name, execute, message string) (*model.Task, error) {
	hash := uuid.New().String()
	task := model.Task{
		Name:    name,
		Execute: execute,
		Message: message,
		Hash:    hash,
		Active:  true,
	}
	if err := s.repo.Create(&task); err != nil {
		return nil, fmt.Errorf("failed to create task in DB: %+v", err)
	}

	if _, err := s.tm.RegisterTask(hash, name, execute, message); err != nil {
		if deleteErr := s.repo.Delete(fmt.Sprint(task.ID)); deleteErr != nil {
			logger.Warnf("failed to rollback task %d after cron registration error: %v", task.ID, deleteErr)
		}
		return nil, fmt.Errorf("failed to register task in cron: %+v", err)
	}

	return &task, nil
}

func (s *TaskService) GetTasks(params httpx.QueryParams) ([]model.Task, int64, error) {
	return s.repo.GetAll(params)
}

func (s *TaskService) GetTaskByID(id string) (*model.Task, error) {
	return s.repo.GetByID(id)
}

func (s *TaskService) GetTaskLogs(hash string) []LogEntry {
	return s.tm.GetLogs(hash)
}

func (s *TaskService) SetTaskActiveStatus(id string, active bool) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("task with id %s not found: %v", id, err)
	}

	if task.Active == active {
		logger.Infof("task with id %s is already in desired state", id)
		return nil
	}

	if active {
		if _, err := s.tm.RegisterTask(task.Hash, task.Name, task.Execute, task.Message); err != nil {
			return fmt.Errorf("failed to register task in cron: %v", err)
		}
	} else {
		if err := s.tm.RemoveTaskFromCronByHash(task.Hash); err != nil {
			logger.Warnf("failed to remove task: %v", err)
		}
	}

	task.Active = active
	if err := s.repo.Update(task); err != nil {
		return err
	}

	statusStr := "ACTIVE"
	if !active {
		statusStr = "INACTIVE"
	}
	s.tm.AddLog(task.Hash, fmt.Sprintf("System: Task is now %s", statusStr))
	return nil
}

func (s *TaskService) UpdateTask(id string, name, execute, message string, active bool) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("task with id %s not found: %v", id, err)
	}

	if task.Name == name && task.Execute == execute && task.Message == message && task.Active == active {
		return nil
	}

	wasActive := task.Active
	if wasActive {
		if err := s.tm.RemoveTaskFromCronByHash(task.Hash); err != nil {
			logger.Warnf("failed to remove task: %v", err)
		}
	}

	task.Name = name
	task.Execute = execute
	task.Message = message
	task.Active = active
	task.Code = ""
	if err := s.repo.Update(task); err != nil {
		if wasActive {
			if _, registerErr := s.tm.RegisterTask(task.Hash, task.Name, task.Execute, task.Message); registerErr != nil {
				logger.Warnf("failed to restore cron task after update error: %v", registerErr)
			}
		}
		return fmt.Errorf("failed to update task in DB: %v", err)
	}

	if task.Active {
		if _, err := s.tm.RegisterTask(task.Hash, name, execute, message); err != nil {
			return fmt.Errorf("failed to re-register updated task in cron: %v", err)
		}
	}

	statusStr := "ACTIVE"
	if !active {
		statusStr = "INACTIVE"
	}
	s.tm.AddLog(task.Hash, fmt.Sprintf("System: Task is now %s", statusStr))
	logger.Infof("updated task: %s", name)
	return nil
}

func (s *TaskService) DeleteTask(id string) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("task with id %s not found: %v", id, err)
	}

	if task.Active {
		if err := s.tm.RemoveTaskFromCronByHash(task.Hash); err != nil {
			logger.Warnf("failed to remove task: %v", err)
		}
	}

	return s.repo.Delete(id)
}

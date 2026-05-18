package repository

import (
	"gorm.io/gorm"

	"github.com/HugManh/cronjob/internal/model"
	"github.com/HugManh/cronjob/pkg/httpx"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db}
}

func (r *TaskRepository) Create(task *model.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepository) GetAll(params httpx.QueryParams) ([]model.Task, int64, error) {
	var tasks []model.Task
	var total int64
	offset := (params.Page - 1) * params.Limit

	r.db.Model(&model.Task{}).Count(&total)

	query := r.db.Model(&model.Task{}).Offset(offset).Limit(params.Limit)
	if params.Sort != "" {
		query = query.Order(params.Sort)
	}

	err := query.Find(&tasks).Error
	return tasks, total, err
}

func (r *TaskRepository) GetByID(id string) (*model.Task, error) {
	var task model.Task
	if err := r.db.Where("id = ?", id).First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) GetByActive(isActive bool) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Where("active = ?", isActive).Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) Update(task *model.Task) error {
	return r.db.Model(&model.Task{}).
		Where("id = ?", task.ID).
		Select("Name", "Execute", "Message", "Hash", "Active", "Code").
		Updates(task).Error
}

func (r *TaskRepository) Delete(id string) error {
	return r.db.Delete(&model.Task{}, id).Error
}

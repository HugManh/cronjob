package repository

import (

	"github.com/HugManh/cronjob/internal/common/request"
	"github.com/HugManh/cronjob/internal/tasks/model"
	"gorm.io/gorm"
)

type TaskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) *TaskRepo {
	return &TaskRepo{db}
}

func (r *TaskRepo) Create(task *model.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepo) GetAll(params request.QueryParams) ([]model.Task, int64, error) {
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

func (r *TaskRepo) GetByID(id string) (*model.Task, error) {
	var task model.Task
	if err := r.db.Where("id = ?", id).First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepo) GetTaskByActive(isActive bool) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Where("active = ?", true).Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepo) Update(task *model.Task) error {
	return r.db.Model(&model.Task{}).
		Where("id = ?", task.ID).
		Select("Name", "Schedule", "Message", "Hash", "Active").
		Updates(task).Error
}

func (r *TaskRepo) Delete(id string) error {
	return r.db.Delete(&model.Task{}, id).Error
}

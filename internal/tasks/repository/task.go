package repository

import (
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

func (r *TaskRepo) GetAll() ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

// func (r *TaskRepo) GetByID(id uint) (*model.Task, error) {
// 	var task model.Task
// 	if err := r.db.First(&task, id).Error; err != nil {
// 		return nil, err
// 	}
// 	return &task, nil
// }

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
	return r.db.Save(task).Error
}

func (r *TaskRepo) Delete(id string) error {
	return r.db.Delete(&model.Task{}, id).Error
}

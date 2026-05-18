package repository

import (
	"gorm.io/gorm"

	"github.com/HugManh/cronjob/internal/model"
	"github.com/HugManh/cronjob/pkg/httpx"
)

type SlackRepository struct {
	db *gorm.DB
}

func NewSlackRepository(db *gorm.DB) *SlackRepository {
	return &SlackRepository{db}
}

func (r *SlackRepository) Create(slack *model.Slack) error {
	return r.db.Create(slack).Error
}

func (r *SlackRepository) GetAll(params httpx.QueryParams) ([]model.Slack, int64, error) {
	var slacks []model.Slack
	var total int64
	offset := (params.Page - 1) * params.Limit

	r.db.Model(&model.Slack{}).Count(&total)

	query := r.db.Model(&model.Slack{}).Offset(offset).Limit(params.Limit)
	if params.Sort != "" {
		query = query.Order(params.Sort)
	}

	err := query.Find(&slacks).Error
	return slacks, total, err
}

func (r *SlackRepository) GetByID(id string) (*model.Slack, error) {
	var slack model.Slack
	if err := r.db.Where("id = ?", id).First(&slack).Error; err != nil {
		return nil, err
	}
	return &slack, nil
}

func (r *SlackRepository) Delete(id string) error {
	return r.db.Delete(&model.Slack{}, id).Error
}

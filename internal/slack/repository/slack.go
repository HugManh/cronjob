package repository

import (
	"github.com/HugManh/cronjob/internal/common/request"
	"github.com/HugManh/cronjob/internal/slack/model"
	"gorm.io/gorm"
)

type SlackRepo struct {
	db *gorm.DB
}

func NewSlackRepo(db *gorm.DB) *SlackRepo {
	return &SlackRepo{db}
}

func (r *SlackRepo) Create(Slack *model.Slack) error {
	return r.db.Create(Slack).Error
}

func (r *SlackRepo) GetAll(params request.QueryParams) ([]model.Slack, int64, error) {
	var Slacks []model.Slack
	var total int64
	offset := (params.Page - 1) * params.Limit

	r.db.Model(&model.Slack{}).Count(&total)

	query := r.db.Model(&model.Slack{}).Offset(offset).Limit(params.Limit)
	if params.Sort != "" {
		query = query.Order(params.Sort)
	}

	err := query.Find(&Slacks).Error
	return Slacks, total, err
}

func (r *SlackRepo) GetByID(id string) (*model.Slack, error) {
	var Slack model.Slack
	if err := r.db.Where("id = ?", id).First(&Slack).Error; err != nil {
		return nil, err
	}
	return &Slack, nil
}

func (r *SlackRepo) Delete(id string) error {
	return r.db.Delete(&model.Slack{}, id).Error
}

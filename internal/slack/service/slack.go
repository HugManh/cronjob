package service

import (
	"github.com/HugManh/cronjob/internal/common/request"
	"github.com/HugManh/cronjob/internal/slack/dto"
	"github.com/HugManh/cronjob/internal/slack/model"
	"github.com/HugManh/cronjob/internal/slack/repository"
)

type SlackService struct {
	repo *repository.SlackRepo
}

func NewService(r *repository.SlackRepo) *SlackService {
	return &SlackService{repo: r}
}

// AddSlack adds a new task to the cron scheduler and saves it to the database
func (s *SlackService) CreateSlack(req dto.CreateSlackRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	slack := &model.Slack{
		BotToken: req.BotToken,
		ChatID:   req.ChatID,
	}

	if err := s.repo.Create(slack); err != nil {
		return err
	}

	return nil
}

func (s *SlackService) GetSlacks(params request.QueryParams) ([]model.Slack, int64, error) {
	return s.repo.GetAll(params)
}

func (s *SlackService) GetSlackById(id string) (*model.Slack, error) {
	return s.repo.GetByID(id)
}

package service

import (
	"github.com/HugManh/cronjob/internal/common/request"
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
func (s *SlackService) CreateSlack(bot_token, chat_id string) error {
	slack := &model.Slack{
		BotToken: bot_token,
		ChatID:   chat_id,
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

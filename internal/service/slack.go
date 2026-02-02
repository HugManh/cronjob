package service

import (
	dto_slack "github.com/HugManh/cronjob/internal/dto/slack"
	"github.com/HugManh/cronjob/internal/model"
	"github.com/HugManh/cronjob/internal/repository"
	"github.com/HugManh/cronjob/pkg/https"
)

type SlackService struct {
	repo *repository.SlackRepo
}

func NewService1(r *repository.SlackRepo) *SlackService {
	return &SlackService{repo: r}
}

// AddSlack adds a new task to the cron scheduler and saves it to the database
func (s *SlackService) CreateSlack(req dto_slack.CreateSlackRequest) error {
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

func (s *SlackService) GetSlacks(params https.QueryParams) ([]model.Slack, int64, error) {
	return s.repo.GetAll(params)
}

func (s *SlackService) GetSlackById(id string) (*model.Slack, error) {
	return s.repo.GetByID(id)
}

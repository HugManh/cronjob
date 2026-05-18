package service

import (
	slackdto "github.com/HugManh/cronjob/internal/dto/slack"
	"github.com/HugManh/cronjob/internal/model"
	"github.com/HugManh/cronjob/internal/repository"
	"github.com/HugManh/cronjob/pkg/httpx"
)

type SlackService struct {
	repo *repository.SlackRepository
}

func NewSlackService(r *repository.SlackRepository) *SlackService {
	return &SlackService{repo: r}
}

// CreateSlack validates and stores a Slack configuration.
func (s *SlackService) CreateSlack(req slackdto.CreateSlackRequest) error {
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

func (s *SlackService) GetSlacks(params httpx.QueryParams) ([]model.Slack, int64, error) {
	return s.repo.GetAll(params)
}

func (s *SlackService) GetSlackByID(id string) (*model.Slack, error) {
	return s.repo.GetByID(id)
}

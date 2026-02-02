package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/slack-go/slack"
)

type ISlack interface {
	SendMessage(channelID, message string) error
}

type slackClient struct {
	client *slack.Client
}

func New(token string) ISlack {
	if token == "" {
		return nil
	}
	return &slackClient{
		client: slack.New(token),
	}
}

func (c *slackClient) SendMessage(channelID, message string) error {
	// Gửi tin nhắn
	channel, timestamp, err := c.client.PostMessage(
		channelID,
		slack.MsgOptionText(message, false),
	)
	if err != nil {
		return err
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("[%s] Channel: %s Timestamp: %s", now, channel, timestamp)
	return nil
}

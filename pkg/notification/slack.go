package main

import (
	"time"

    log "github.com/sirupsen/logrus"

	"github.com/slack-go/slack"
)

func newSlackClient(botToken string) *slack.Client {
	if botToken != "" {
		return slack.New(botToken)
	}
	log.Fatal("Missing Slack Bot Token variable")
	return nil
}

func sendMessage(api *slack.Client, channelID, message string) {
	// Gửi tin nhắn
	channel, timestamp, err := api.PostMessage(
		channelID,
		slack.MsgOptionText(message, false),
	)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("[%s] Channel: %s Timestamp: %s", now, channel, timestamp)
}

package dto

import (
	"fmt"
)

type AddSlackRequest struct {
	BotToken string `json:"bot_token"`
	ChatID   string `json:"chat_id"`
}

func (r *AddSlackRequest) Validate() error {
	if r.BotToken == "" {
		return fmt.Errorf("bot_token is required")
	}
	if r.ChatID == "" {
		return fmt.Errorf("chat_id is required")
	}
	return nil
}

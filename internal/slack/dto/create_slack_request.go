package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateSlackRequest struct {
	BotToken string `json:"bot_token"`
	ChatID   string `json:"chat_id"`
}

func (r CreateSlackRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.BotToken, validation.Required),
		validation.Field(&r.ChatID, validation.Required),
	)
}

package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateSlackRequest struct {
	BotToken string `json:"botToken"`
	ChatID   string `json:"chatID"`
}

func (r CreateSlackRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.BotToken, validation.Required),
		validation.Field(&r.ChatID, validation.Required),
	)
}

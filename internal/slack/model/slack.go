package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedAt time.Time      `gorm:"column:created_at; autoCreateTime; <-:create" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at; autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at; index" json:"-"`
}

type Slack struct {
	ID       int    `gorm:"column:id; primary_key; not null" json:"id"`
	BotToken string `gorm:"column:bot_token; type:varchar(255); not null; index" json:"botToken"`
	ChatID   string `gorm:"column:chat_id; type:varchar(255); not null" json:"chatID"`
	BaseModel
}

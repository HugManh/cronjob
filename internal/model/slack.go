package model

type Slack struct {
	ID       int    `gorm:"column:id; primary_key; not null" json:"id"`
	BotToken string `gorm:"column:bot_token; type:varchar(255); not null; index" json:"botToken"`
	ChatID   string `gorm:"column:chat_id; type:varchar(255); not null" json:"chatID"`
	BaseModel
}

package model

type Task struct {
	ID      int    `gorm:"column:id; primary_key; not null" json:"id"`
	Name    string `gorm:"column:name; type:varchar(255); not null; index" json:"name"`
	Execute string `gorm:"column:execute; type:varchar(255); not null" json:"execute"`
	Message string `gorm:"column:message; type:text" json:"message"`
	Hash    string `gorm:"column:hash; type:varchar(64); uniqueIndex; not null" json:"hash"`
	Active  bool   `gorm:"column:active; not null; default:true" json:"active"`
    Code    string `gorm:"column:code; type:varchar(64); null" json:"code"`
	BaseModel
}

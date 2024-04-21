package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Id        uuid.UUID `json:"id" gorm:"primaryKey"`
	Text      string    `json:"text" gorm:"column:text"`
	Status    string    `json:"status" gorm:"column:status"`
	SentBy    uuid.UUID `json:"sent_by" gorm:"column:sent_by"`
	SentTo    uuid.UUID `json:"sent_to" gorm:"column:sent_to"`
	RepliedOn string    `json:"replied_on" gorm:"column:replied_on"`
	ChatId    string    `json:"chat_id" gorm:"column:chat_id"`
	Chat      Chat      `gorm:"foreignKey:ChatId"`
}

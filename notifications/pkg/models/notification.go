package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	Id   uuid.UUID `json:"id" gorm:"primaryKey"`
	Text string    `json:"text" gorm:"column:text"`
}

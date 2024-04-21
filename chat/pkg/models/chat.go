package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	Id    uuid.UUID `json:"id" gorm:"primaryKey"`
	User1 uuid.UUID `json:"user_1" gorm:"type:uuid; column:user_1"`
	User2 uuid.UUID `json:"user_2" gorm:"type:uuid; column:user_2"`
}

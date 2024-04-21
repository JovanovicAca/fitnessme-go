package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Exercise struct {
	gorm.Model
	Id          uuid.UUID `json:"id" gorm:"primaryKey"`
	Description string    `json:"description" gorm:"column:description"`
	CreatedBy   uuid.UUID `json:"createdBy" gorm:"column:created_by"`
	Name        string    `json:"name" gorm:"column:name"`
}

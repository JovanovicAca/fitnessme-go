package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExerciseGroup struct {
	gorm.Model
	Id          uuid.UUID `json:"id" gorm:"primaryKey"`
	GroupName   string    `json:"group_name" gorm:"column:group_name"`
	Description string    `json:"description" gorm:"column:description"`
}

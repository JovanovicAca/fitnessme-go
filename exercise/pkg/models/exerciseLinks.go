package models

import (
	"gorm.io/gorm"
)

type ExerciseLinks struct {
	gorm.Model
	ExerciseID string `json:"exercise_id" gorm:"column:exercise_id"`
	Link       string `json:"link" gorm:"column:link"`
}

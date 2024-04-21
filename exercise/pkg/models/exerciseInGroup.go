package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExerciseInGroup struct {
	gorm.Model
	Id              uuid.UUID     `json:"id" gorm:"primaryKey"`
	ExerciseID      uuid.UUID     `json:"exercise_id" gorm:"type:uuid; column:exercise_id"`
	ExerciseGroupID uuid.UUID     `json:"exercise_group_id" gorm:"type:uuid; column:exercise_group_id"`
	SequenceOrder   int           `json:"sequence_order" gorm:"column:sequence_order"`
	Exercise        Exercise      `gorm:"foreignKey:ExerciseID"`
	ExerciseGroup   ExerciseGroup `gorm:"foreignKey:ExerciseGroupID"`
}

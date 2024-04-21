package dto

import (
	"github.com/google/uuid"
)

type ExerciseEditDTO struct {
	Name          string    `json:"name"`
	ExerciseGroup string    `json:"exercise_group"`
	Description   string    `json:"description"`
	CreatedBy     uuid.UUID `json:"created_by"`
	SequenceOrder int       `json:"sequence_order"`
	Link          string    `json:"link"`
}

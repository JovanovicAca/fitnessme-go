package dto

import (
	"github.com/google/uuid"
)

type ExerciseInGroupDTO struct {
	Exercise      uuid.UUID `json:"exercise"`
	ExerciseGroup uuid.UUID `json:"exercise_group"`
	SequenceOrder int       `json:"sequence_order"`
}

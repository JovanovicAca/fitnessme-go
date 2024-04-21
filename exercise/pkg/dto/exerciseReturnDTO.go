package dto

import (
	"github.com/google/uuid"
)

type ExerciseReturnDTO struct {
	Id            uuid.UUID
	Name          string
	ExerciseGroup string
	Description   string
	CreatedBy     uuid.UUID
	SequenceOrder int
	Link          string
}

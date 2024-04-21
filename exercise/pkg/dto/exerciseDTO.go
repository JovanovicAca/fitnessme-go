package dto

import (
	"fitnessme/exercise/pkg/models"

	"github.com/google/uuid"
)

type ExerciseDTO struct {
	Name          string    `json:"name"`
	ExerciseGroup uuid.UUID `json:"exercise_group"`
	Description   string    `json:"description"`
	CreatedBy     uuid.UUID `json:"created_by"`
	SequenceOrder int       `json:"sequence_order"`
	Link          string    `json:"link"`
}

func (dto *ExerciseDTO) ToExerciseModel() models.Exercise {
	return models.Exercise{
		Id:          uuid.New(),
		Description: dto.Description,
		CreatedBy:   dto.CreatedBy,
		Name:        dto.Name,
	}
}

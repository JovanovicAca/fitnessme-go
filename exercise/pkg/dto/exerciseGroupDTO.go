package dto

import (
	"github.com/google/uuid"
)

type ExerciseGroupDTO struct {
	Id          uuid.UUID `json:"id"`
	GroupName   string    `json:"group_name"`
	Description string    `json:"description"`
}

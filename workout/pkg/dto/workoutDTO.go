package dto

import "github.com/google/uuid"

type WorkoutDTO struct {
	UserId     uuid.UUID `json:"userId"`
	ExerciseId uuid.UUID `json:"exerciseId"`
	WoroutDate string    `json:"workoutDate"`
	Sets       int       `json:"sets"`
	Reps       int       `json:"reps"`
	Weight     float32   `json:"weight"`
	Duration   float32   `json:"duration"`
}

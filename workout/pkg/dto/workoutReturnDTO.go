package dto

import (
	"time"

	"github.com/google/uuid"
)

type WorkoutReturnDTO struct {
	UserId      uuid.UUID
	Exercise    string
	WorkoutDate time.Time
	Sets        int
	Reps        int
	Weight      float32
	Duration    float32
}

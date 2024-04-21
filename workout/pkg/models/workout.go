package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Workout struct {
	gorm.Model
	Id          uuid.UUID `json:"id" gorm:"primaryKey"`
	WorkoutId   uuid.UUID `json:"workout_id" gorm:"column:workout_id"`
	UserId      uuid.UUID `json:"user_id" gorm:"column:user_id"`
	ExerciseId  uuid.UUID `json:"exercise_id" gorm:"column:exercise_id"`
	WorkoutDate time.Time `json:"workout_date" gorm:"column:workout_date"`
	Sets        int       `json:"sets" gorm:"column:sets"`
	Reps        int       `json:"reps" gorm:"column:reps"`
	Weight      float32   `json:"weight" gorm:"column:weight"`
	Duration    float32   `json:"duration" gorm:"column:duration"`
}

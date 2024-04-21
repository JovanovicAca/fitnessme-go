package repository

import (
	"fitnessme/workout/pkg/db"
	"fitnessme/workout/pkg/models"

	"github.com/pkg/errors"
)

type WorkoutRepository interface {
	Create(workout models.Workout) error
	GetWorkoutsForUser(user_id string) ([]models.Workout, error)
	DeleteWorkout(id string) error
}

type workoutRepository struct{ handler db.Handler }

func NewWorkoutRepository(handler db.Handler) *workoutRepository {
	return &workoutRepository{handler: handler}
}

func (w *workoutRepository) DeleteWorkout(id string) error {
	if err := w.handler.DB.Where("exercise_id = ?", id).Delete(&models.Workout{}).Error; err != nil {
		return err
	}
	return nil
}

func (w *workoutRepository) GetWorkoutsForUser(user_id string) ([]models.Workout, error) {
	var workouts []models.Workout
	result := w.handler.DB.Where("user_id = ?", user_id).Find(&workouts)
	if result.Error != nil {
		return nil, result.Error
	}
	return workouts, nil
}

func (w *workoutRepository) Create(workout models.Workout) error {
	if err := w.handler.DB.Create(&workout).Error; err != nil {
		return errors.Wrap(err, "failed to save workout")
	}
	return nil
}

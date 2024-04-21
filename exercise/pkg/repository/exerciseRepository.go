package repository

import (
	"fitnessme/exercise/pkg/db"
	"fitnessme/exercise/pkg/dto"
	"fitnessme/exercise/pkg/models"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ExerciseRepository interface {
	FindInGroupById(id uuid.UUID) (models.ExerciseInGroup, error)
	FindGroupById(id uuid.UUID) (models.ExerciseGroup, error)
	SaveExercise(exercise models.Exercise) error
	SaveExerciseInGroup(exerciseInGroup models.ExerciseInGroup) error
	SaveExerciseLink(exerciseLink models.ExerciseLinks) error
	SaveExerciseGroup(exerciseGroup models.ExerciseGroup) error
	GetAllGroups() ([]models.ExerciseGroup, error)
	GetAllExerciseIdsByGroupId(id uuid.UUID) ([]uuid.UUID, error)
	FindExerciseById(id uuid.UUID) (models.Exercise, error)
	GetGroupNameById(id uuid.UUID) string
	FindSequenceOrderForExercise(group_id, exercise_id uuid.UUID) int
	FindLinkByExerciseId(exercise_id uuid.UUID) string
	FindAllExercisesByGroup(groupId string) ([]dto.ExerciseReturnDTO, error)
	CheckIfExerciseExist(group_id uuid.UUID, exercise_name string) bool
	GetGroupByName(group_name string) (models.ExerciseGroup, error)
	GetExerciseNameById(exercise_id string) (string, error)
	GetAllExercises() ([]models.Exercise, error)
	DeleteExerciseWithAssociations(id string) error
	UpdateExercise(id string, exercise dto.ExerciseEditDTO) error
}

type exerciseRepository struct{ handler db.Handler }

func NewExerciseRepository(handler db.Handler) ExerciseRepository {
	return &exerciseRepository{handler: handler}
}

func (e *exerciseRepository) UpdateExercise(id string, exercise dto.ExerciseEditDTO) error {
	tx := e.handler.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Model(&models.Exercise{}).Where("id = ?", id).Updates(models.Exercise{
		Name:        exercise.Name,
		Description: exercise.Description,
		CreatedBy:   exercise.CreatedBy,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	var exerciseGroup models.ExerciseGroup
	if err := tx.Where("group_name = ?", exercise.ExerciseGroup).First(&exerciseGroup).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&models.ExerciseInGroup{}).Where("exercise_id = ?", id).Updates(map[string]interface{}{
		"sequence_order":    exercise.SequenceOrder,
		"exercise_group_id": exerciseGroup.Id,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&models.ExerciseLinks{}).Where("exercise_id = ?", id).Update("link", exercise.Link).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (e *exerciseRepository) DeleteExerciseWithAssociations(id string) error {
	tx := e.handler.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("exercise_id = ?", id).Delete(&models.ExerciseInGroup{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("exercise_id = ?", id).Delete(&models.ExerciseLinks{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&models.Exercise{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (e *exerciseRepository) GetExerciseNameById(exercise_id string) (string, error) {
	var exercise models.Exercise
	if err := e.handler.DB.Where("id = ?", exercise_id).First(&exercise).Error; err != nil {
		return exercise.Name, errors.Wrap(err, "exercise not found")
	}
	return exercise.Name, nil
}

func (e *exerciseRepository) GetGroupByName(group_name string) (models.ExerciseGroup, error) {
	var group models.ExerciseGroup
	if err := e.handler.DB.Where("LOWER(group_name) = LOWER(?)", group_name).First(&group).Error; err != nil {
		return group, errors.Wrap(err, "group not found")
	}
	return group, nil
}

func (e *exerciseRepository) CheckIfExerciseExist(group_id uuid.UUID, exercise_name string) bool {
	var exercise models.Exercise
	if err := e.handler.DB.Where("name = ?", exercise_name).First(&exercise).Error; err != nil {
		return false
	}

	var exerciseInGroup models.ExerciseInGroup
	result := e.handler.DB.Where("exercise_id = ? AND exercise_group_id = ?", exercise.Id, group_id).First(&exerciseInGroup)
	return result.Error == nil
}

func (e *exerciseRepository) FindAllExercisesByGroup(groupId string) ([]dto.ExerciseReturnDTO, error) {
	var exercises []dto.ExerciseReturnDTO
	if err := e.handler.DB.
		Table("exercise_in_groups").
		Select("exercises.id, exercises.name, exercise_groups.group_name as exercise_group, exercises.description, exercises.created_by, exercise_in_groups.sequence_order, exercise_links.link").
		Joins("join exercise_groups on exercise_groups.id = exercise_in_groups.exercise_group_id").
		Joins("join exercises on exercises.id = exercise_in_groups.exercise_id").
		Joins("join exercise_links on exercises.id = exercise_links.exercise_id").
		Where("exercise_in_groups.exercise_group_id = ? AND exercises.deleted_at IS NULL", groupId).
		Scan(&exercises).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get exercises")
	}

	return exercises, nil
}

func (e *exerciseRepository) FindLinkByExerciseId(exercise_id uuid.UUID) string {
	var link string
	if err := e.handler.DB.Model(&models.ExerciseLinks{}).Where("exercise_id = ?", exercise_id).Pluck("link", &link).Error; err != nil {
		return ""
	}
	return link
}

func (e *exerciseRepository) FindSequenceOrderForExercise(group_id, exercise_id uuid.UUID) int {
	var sequenceOrder int
	if err := e.handler.DB.Model(&models.ExerciseInGroup{}).Where("exercise_id = ? AND exercise_group_id", exercise_id, group_id).Pluck("sequence_order", &sequenceOrder).Error; err != nil {
		return 0
	}
	return sequenceOrder
}

func (e *exerciseRepository) GetGroupNameById(id uuid.UUID) string {
	var groupName string
	if err := e.handler.DB.Model(&models.ExerciseGroup{}).Where("id = ?", id).Pluck("group_name", &groupName).Error; err != nil {
		return ""
	}
	return groupName
}

func (e *exerciseRepository) FindExerciseById(id uuid.UUID) (models.Exercise, error) {
	var exercise models.Exercise
	if err := e.handler.DB.Where("id = ?", id).First(&exercise).Error; err != nil {
		return exercise, errors.Wrap(err, "exercise not found")
	}
	return exercise, nil
}

func (e *exerciseRepository) GetAllExerciseIdsByGroupId(id uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	if err := e.handler.DB.Where("exercise_group_id = ?", id).Pluck("exercise_id", &ids).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get exercise ids")
	}
	return ids, nil
}

func (e *exerciseRepository) GetAllExercises() ([]models.Exercise, error) {
	var exercises []models.Exercise
	if err := e.handler.DB.Find(&exercises).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get all exercises")
	}
	return exercises, nil
}

func (e *exerciseRepository) GetAllGroups() ([]models.ExerciseGroup, error) {
	var exercise_groups []models.ExerciseGroup
	if err := e.handler.DB.Find(&exercise_groups).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get all groups")
	}
	return exercise_groups, nil
}

func (e *exerciseRepository) SaveExerciseGroup(exerciseGroup models.ExerciseGroup) error {
	if err := e.handler.DB.Create(&exerciseGroup).Error; err != nil {
		return errors.Wrap(err, "failed to save exercise group")
	}
	return nil
}

func (e *exerciseRepository) SaveExerciseLink(exerciseLink models.ExerciseLinks) error {
	if err := e.handler.DB.Create(&exerciseLink).Error; err != nil {
		return errors.Wrap(err, "failed to save exercise link")
	}
	return nil
}

func (e *exerciseRepository) SaveExerciseInGroup(exerciseInGroup models.ExerciseInGroup) error {
	if err := e.handler.DB.Create(&exerciseInGroup).Error; err != nil {
		return errors.Wrap(err, "failed to save exercise in group")
	}
	return nil
}

func (e *exerciseRepository) SaveExercise(exercise models.Exercise) error {
	if err := e.handler.DB.Create(&exercise).Error; err != nil {
		return errors.Wrap(err, "failed to save exercise")
	}
	return nil
}

func (e *exerciseRepository) FindInGroupById(id uuid.UUID) (models.ExerciseInGroup, error) {
	var exerciseInGroup models.ExerciseInGroup
	if err := e.handler.DB.Where("id = ?", id).First(&exerciseInGroup).Error; err != nil {
		return exerciseInGroup, errors.Wrap(err, "exercise in group not found")
	}
	return exerciseInGroup, nil
}

func (e *exerciseRepository) FindGroupById(id uuid.UUID) (models.ExerciseGroup, error) {
	var exerciseGroup models.ExerciseGroup
	if err := e.handler.DB.Where("id = ?", id).First(&exerciseGroup).Error; err != nil {
		return exerciseGroup, errors.Wrap(err, "exercise group not found")
	}
	return exerciseGroup, nil
}

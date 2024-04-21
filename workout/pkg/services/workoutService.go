package services

import (
	"encoding/json"
	"fitnessme/utils"
	"fitnessme/workout/pkg/dto"
	"fitnessme/workout/pkg/models"
	"fitnessme/workout/pkg/repository"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type WorkoutService interface {
	Ping(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	GetWorkouts(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type workoutService struct {
	repo       repository.WorkoutRepository
	jwtWrapper *utils.JwtWrapper
}

func NewWorkoutService(repo repository.WorkoutRepository, jwtWrapper *utils.JwtWrapper) WorkoutService {
	return &workoutService{repo: repo, jwtWrapper: jwtWrapper}
}

func (wo *workoutService) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id query parameter", http.StatusBadRequest)
		return
	}

	err := wo.repo.DeleteWorkout(id)
	if err != nil {
		http.Error(w, "Error in deleting", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (wo *workoutService) GetWorkouts(w http.ResponseWriter, r *http.Request) {
	var idMessage string
	var valid bool

	valid, idMessage = wo.getUserId(r)

	if !valid {
		http.Error(w, idMessage, http.StatusUnauthorized)
		return
	}

	workouts, err := wo.repo.GetWorkoutsForUser(idMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	groupedWorkoutsMap := make(map[string][]dto.WorkoutReturnDTO)

	exerciseServiceURL := "http://localhost:3001/exercise?id="

	for _, workout := range workouts {
		exerciseURL := fmt.Sprintf("%s%s", exerciseServiceURL, workout.ExerciseId)

		resp, err := http.Get(exerciseURL)
		if err != nil {
			http.Error(w, "Failed to get exercise name", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var exercise struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&exercise); err != nil {
			http.Error(w, "Failed to decode exercise details", http.StatusInternalServerError)
			return
		}

		workoutDTO := dto.WorkoutReturnDTO{
			UserId:      workout.UserId,
			Exercise:    exercise.Name,
			WorkoutDate: workout.WorkoutDate,
			Sets:        workout.Sets,
			Reps:        workout.Reps,
			Weight:      workout.Weight,
			Duration:    workout.Duration,
		}

		groupedWorkoutsMap[workout.WorkoutId.String()] = append(groupedWorkoutsMap[workout.WorkoutId.String()], workoutDTO)
	}

	var groupedWorkouts []dto.GroupedWorkouts
	for workoutID, exercises := range groupedWorkoutsMap {
		groupedWorkout := dto.GroupedWorkouts{
			WorkoutID: workoutID,
			Exercises: exercises,
		}
		groupedWorkouts = append(groupedWorkouts, groupedWorkout)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(groupedWorkouts)
}

func (wo *workoutService) Create(w http.ResponseWriter, r *http.Request) {
	var workoutDTOs []dto.WorkoutDTO
	err1 := json.NewDecoder(r.Body).Decode(&workoutDTOs)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	workout_id := uuid.New()
	for _, workoutDTO := range workoutDTOs {
		workoutDate, err := time.Parse("2006-01-02", workoutDTO.WoroutDate)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}

		workout := models.Workout{
			Id:          uuid.New(),
			WorkoutId:   workout_id,
			UserId:      workoutDTO.UserId,
			ExerciseId:  workoutDTO.ExerciseId,
			WorkoutDate: workoutDate,
			Sets:        workoutDTO.Sets,
			Reps:        workoutDTO.Reps,
			Weight:      workoutDTO.Weight,
			Duration:    workoutDTO.Duration,
		}

		err = wo.repo.Create(workout)
		if err != nil {
			http.Error(w, "Failed to save workout", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Workout added")
}

func (wo *workoutService) Ping(w http.ResponseWriter, r *http.Request) {
	print("success")
}

func (wo *workoutService) getUserId(r *http.Request) (bool, string) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false, "Authorization header is required"
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return false, "Invalid Authorization header format"
	}

	token := headerParts[1]
	claims, err := wo.jwtWrapper.ValidateToken(token)
	if err != nil {
		return false, "Invalid token: " + err.Error()
	}

	return true, claims.Id.String()
}

package services

import (
	"bytes"
	"encoding/json"
	"fitnessme/exercise/pkg/dto"
	"fitnessme/exercise/pkg/models"
	"fitnessme/exercise/pkg/repository"
	"fitnessme/utils"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type ExerciseService interface {
	Ping(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	CreateGroup(w http.ResponseWriter, r *http.Request)
	GetAllGroups(w http.ResponseWriter, r *http.Request)
	GetExerciseByGroupId(w http.ResponseWriter, r *http.Request)
	GetNameById(w http.ResponseWriter, r *http.Request)
	GetAllExercises(w http.ResponseWriter, r *http.Request)
	DeleteExercise(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

type exerciseService struct {
	repo       repository.ExerciseRepository
	jwtWrapper *utils.JwtWrapper
}

func NewExerciseService(repo repository.ExerciseRepository, jwtWrapper *utils.JwtWrapper) ExerciseService {
	return &exerciseService{repo: repo, jwtWrapper: jwtWrapper}
}
func (e *exerciseService) Update(w http.ResponseWriter, r *http.Request) {
	if admin, message := e.checkIfRoleIsAdmin(r); !admin {
		http.Error(w, message, http.StatusUnauthorized)
		return
	}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	bodyString := string(bodyBytes)
	fmt.Println("Received request body:", bodyString)

	reader := bytes.NewReader(bodyBytes)

	var exercise dto.ExerciseEditDTO
	if err := json.NewDecoder(reader).Decode(&exercise); err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id query parameter", http.StatusBadRequest)
		return
	}
	err = e.repo.UpdateExercise(id, exercise)
	if err != nil {
		http.Error(w, "Error in updating", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (e *exerciseService) DeleteExercise(w http.ResponseWriter, r *http.Request) {
	if admin, message := e.checkIfRoleIsAdmin(r); !admin {
		http.Error(w, message, http.StatusUnauthorized)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id query parameter", http.StatusBadRequest)
		return
	}

	err := e.repo.DeleteExerciseWithAssociations(id)
	if err != nil {
		http.Error(w, "Error in deleting", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (e *exerciseService) GetAllExercises(w http.ResponseWriter, r *http.Request) {
	exercises, err := e.repo.GetAllExercises()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	jsonResponse, err := json.Marshal(exercises)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (e *exerciseService) GetNameById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id query parameter", http.StatusBadRequest)
		return
	}
	exerciseName, err := e.repo.GetExerciseNameById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"name": exerciseName}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (e *exerciseService) GetExerciseByGroupId(w http.ResponseWriter, r *http.Request) {
	group_id := r.URL.Query().Get("group")

	exercisesReturn, err := e.repo.FindAllExercisesByGroup(group_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal(exercisesReturn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (e *exerciseService) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := e.repo.GetAllGroups()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(groups)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (e *exerciseService) CreateGroup(w http.ResponseWriter, r *http.Request) {
	if admin, message := e.checkIfRoleIsAdmin(r); !admin {
		http.Error(w, message, http.StatusUnauthorized)
		return
	}

	var exerciseGroup models.ExerciseGroup
	if err := json.NewDecoder(r.Body).Decode(&exerciseGroup); err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	exerciseGroup.Id = uuid.New()

	// check if group already exist
	if _, err := e.repo.GetGroupByName(strings.ToLower(exerciseGroup.GroupName)); err == nil {
		http.Error(w, "Group with that name already exist", http.StatusConflict)
		return
	}

	if err := e.repo.SaveExerciseGroup(exerciseGroup); err != nil {
		http.Error(w, "Failed to save exercise group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Exercise group saved successfully")
}

func (e *exerciseService) Ping(w http.ResponseWriter, r *http.Request) {
	if admin, message := e.checkIfRoleIsAdmin(r); !admin {
		http.Error(w, message, http.StatusUnauthorized)
		return
	}
	print("success")
}

func (e *exerciseService) Create(w http.ResponseWriter, r *http.Request) {
	if admin, message := e.checkIfRoleIsAdmin(r); !admin {
		http.Error(w, message, http.StatusUnauthorized)
		return
	}
	var exerciseDTO dto.ExerciseDTO
	if err := json.NewDecoder(r.Body).Decode(&exerciseDTO); err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	exercise := exerciseDTO.ToExerciseModel()

	//check if exercise exist with that name and group
	if e.repo.CheckIfExerciseExist(exerciseDTO.ExerciseGroup, exercise.Name) {
		http.Error(w, "Exercise with that name already exist in that group", http.StatusBadRequest)
		return
	}

	if err := e.repo.SaveExercise(exercise); err != nil {
		http.Error(w, "Failed to save exercise", http.StatusInternalServerError)
		return
	}

	exerciseInGroup := models.ExerciseInGroup{
		Id:              uuid.New(),
		ExerciseID:      exercise.Id,
		ExerciseGroupID: exerciseDTO.ExerciseGroup,
		SequenceOrder:   exerciseDTO.SequenceOrder,
	}

	if err := e.repo.SaveExerciseInGroup(exerciseInGroup); err != nil {
		http.Error(w, "Failed to save exercise in group", http.StatusInternalServerError)
		return
	}

	if exerciseDTO.Link != "" {
		exerciseLink := models.ExerciseLinks{
			ExerciseID: exercise.Id.String(),
			Link:       exerciseDTO.Link,
		}

		if err := e.repo.SaveExerciseLink(exerciseLink); err != nil {
			http.Error(w, "Failed to save exercise link", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Exercise saved successfully")
}

func (e *exerciseService) checkIfRoleIsAdmin(r *http.Request) (bool, string) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false, "Authorization header is required"
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return false, "Invalid Authorization header format"
	}

	token := headerParts[1]
	claims, err := e.jwtWrapper.ValidateToken(token)
	if err != nil {
		return false, "Invalid token: " + err.Error()
	}

	if claims.Role != "admin" {
		return false, "Unauthorized - Admin role required"
	}

	return true, "Authorized"
}

package services

import (
	"encoding/json"
	"fitnessme/usermanagement/pkg/dto"
	"fitnessme/usermanagement/pkg/models"
	"fitnessme/usermanagement/pkg/repository"
	"fitnessme/utils"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type UserManagementService interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	GetAllAdmins(w http.ResponseWriter, r *http.Request)
}

type userManagementService struct {
	repo       repository.UserRepository
	jwtWrapper *utils.JwtWrapper
}

func NewUserManagementService(repo repository.UserRepository, jwtWrapper *utils.JwtWrapper) UserManagementService {
	return &userManagementService{repo: repo, jwtWrapper: jwtWrapper}
}

func (u *userManagementService) GetAllAdmins(w http.ResponseWriter, r *http.Request) {
	admins, err := u.repo.GetAllAdmins()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(admins)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}

func (u *userManagementService) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var idMessage string
	var valid bool

	valid, idMessage = u.getUserId(r)
	if !valid {
		http.Error(w, idMessage, http.StatusUnauthorized)
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println(updates)
	if dob, exists := updates["date_of_birth"].(string); exists {
		parsedDob, err := time.Parse("2006-01-02", dob)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
		updates["date_of_birth"] = parsedDob
	}

	if err := u.repo.UpdateUser(idMessage, updates); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintln(w, "User updated successfully")

}

func (u *userManagementService) GetUser(w http.ResponseWriter, r *http.Request) {
	var idMessage string
	var valid bool

	valid, idMessage = u.getUserId(r)
	if !valid {
		http.Error(w, idMessage, http.StatusUnauthorized)
	}

	user, err := u.repo.GetUserById(idMessage)
	if err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (u *userManagementService) Register(w http.ResponseWriter, r *http.Request) {
	var userDTO dto.UserDTO
	err := json.NewDecoder(r.Body).Decode(&userDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, errUser := u.repo.FindByEmail(userDTO.Email)
	if errUser == nil {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	_, err = time.Parse("2006-01-02", userDTO.DateOfBirth)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	newUser := userDTO.ToUserModel()

	err = u.repo.Register(newUser)
	if err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User registered successfully")
}

func (u *userManagementService) Login(w http.ResponseWriter, r *http.Request) {
	var loginDTO dto.LoginDTO
	err := json.NewDecoder(r.Body).Decode(&loginDTO)
	if err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	passwordHashed := loginDTO.Password
	var user models.User
	user, err = u.repo.FindByEmail(loginDTO.Email)
	if err != nil {
		http.Error(w, "Failed to find user", http.StatusNotFound)
		return
	}

	password := utils.CheckPasswordHash(passwordHashed, user.Password)
	if !password {
		http.Error(w, "Failed to find user", http.StatusNotFound)
		return
	}

	accessToken, err := u.jwtWrapper.GenerateToken(user)
	if err != nil {
		fmt.Print("2222")
		http.Error(w, "failed to generate access token", http.StatusInternalServerError)
		return
	}
	fmt.Print("111111")
	w.Header().Set("Authorization", "Bearer "+accessToken)
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, "User logged in successfully")
}

func (u *userManagementService) getUserId(r *http.Request) (bool, string) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false, "Authorization header is required"
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return false, "Invalid Authorization header format"
	}

	token := headerParts[1]
	claims, err := u.jwtWrapper.ValidateToken(token)
	if err != nil {
		return false, "Invalid token: " + err.Error()
	}

	return true, claims.Id.String()
}

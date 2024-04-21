package dto

import (
	"fitnessme/usermanagement/pkg/models"
	"fitnessme/utils"
	"time"

	"github.com/google/uuid"
)

type UserDTO struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Address     string `json:"address"`
	DateOfBirth string `json:"date_of_birth"`
	Role        string `json:"role"`
}

func (dto *UserDTO) ToUserModel() models.User {
	dateOfBirth, _ := time.Parse("2006-01-02", dto.DateOfBirth)
	return models.User{
		Id:          uuid.New(),
		Email:       dto.Email,
		Password:    utils.HashPassword(dto.Password),
		Name:        dto.Name,
		Surname:     dto.Surname,
		Address:     dto.Address,
		DateOfBirth: dateOfBirth,
		Blocked:     false,
		Deleted:     false,
		Role:        dto.Role,
	}
}

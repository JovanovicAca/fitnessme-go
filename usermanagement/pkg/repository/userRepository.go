package repository

import (
	"fitnessme/usermanagement/pkg/db"
	"fitnessme/usermanagement/pkg/models"

	"github.com/pkg/errors"
)

type UserRepository interface {
	Register(user models.User) error
	FindByEmail(email string) (models.User, error)
	GetUserById(id string) (models.User, error)
	UpdateUser(id string, updates map[string]interface{}) error
	GetAllAdmins() ([]models.User, error)
}

type userRepository struct{ handler db.Handler }

func NewUserRepository(handler db.Handler) UserRepository {
	return &userRepository{handler: handler}
}

func (u *userRepository) GetAllAdmins() ([]models.User, error) {
	var admins []models.User
	if err := u.handler.DB.Where("role = ?", "admin").Find(&admins).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get all admins")
	}
	return admins, nil
}

func (u *userRepository) UpdateUser(id string, updates map[string]interface{}) error {
	result := u.handler.DB.Model(&models.User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no user found with the given ID")
	}

	return nil
}

func (u *userRepository) GetUserById(id string) (models.User, error) {
	var user models.User
	if err := u.handler.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return user, errors.Wrap(err, "user not found")
	}

	return user, nil
}

func (u *userRepository) Register(user models.User) error {
	if err := u.handler.DB.Create(&user).Error; err != nil {
		return errors.Wrap(err, "failed to save user")
	}
	return nil
}

func (u *userRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	if err := u.handler.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return user, errors.Wrap(err, "user not found")
	}

	return user, nil
}

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id          uuid.UUID `json:"id" gorm:"primaryKey"`
	Email       string    `json:"email" gorm:"column:email"`
	Password    string    `json:"password" gorm:"column:password"`
	Name        string    `json:"name" gorm:"column:name"`
	Surname     string    `json:"surname" gorm:"column:surname"`
	Address     string    `json:"address" gorm:"column:address"`
	DateOfBirth time.Time `json:"date_of_birth" gorm:"column:date_of_birth"`
	Blocked     bool      `json:"blocked" gorm:"column:blocked"`
	Deleted     bool      `json:"deleted" gorm:"column:deleted"`
	Role        string    `json:"role" gorm:"column:role"`
}

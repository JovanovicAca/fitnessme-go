package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"fitnessme/workout/pkg/models"
)

type Handler struct {
	DB *gorm.DB
}

func Init(host, dbname, user, password string, port int) Handler {
	con := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", host, port, dbname, user, password)
	db, err := gorm.Open(postgres.Open(con), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&models.Workout{})
	return Handler{db}
}

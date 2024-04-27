package main

import (
	"context"
	chat "fitnessme/chat/cmd"
	exercise "fitnessme/exercise/cmd"
	notification "fitnessme/notifications/cmd"
	usermanagement "fitnessme/usermanagement/cmd"
	"fitnessme/utils"
	workout "fitnessme/workout/cmd"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("hello world")
	fmt.Println()

	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := os.Getenv("JWT_SECRET_KEY")
	issuer := os.Getenv("JWT_ISSUE")
	expirationHours, _ := strconv.ParseInt(os.Getenv("JWT_EXPIRATION_HOURS"), 10, 64)

	jw := &utils.JwtWrapper{
		SecretKey:       secretKey,
		Issuer:          issuer,
		ExpirationHours: expirationHours,
	}

	userApp := usermanagement.New(jw)
	go func() {
		err := userApp.Start(context.Background())
		if err != nil {
			fmt.Println("failed to start app: ", err)
		}
	}()

	exerciseApp := exercise.New(jw)
	go func() {
		err := exerciseApp.Start(context.Background())
		if err != nil {
			fmt.Println("failed to start app: ", err)
		}
	}()

	workoutApp := workout.New(jw)
	go func() {
		err := workoutApp.Start(context.Background())
		if err != nil {
			fmt.Println("failed to start app: ", err)
		}
	}()

	chatApp := chat.New(jw)
	go func() {
		err := chatApp.Start(context.Background())
		if err != nil {
			fmt.Println("failed to start app: ", err)
		}
	}()

	notificationApp := notification.New()
	go func() {
		err := notificationApp.Start(context.Background())
		if err != nil {
			fmt.Println("failed to start app: ", err)
		}
	}()

	select {}
}

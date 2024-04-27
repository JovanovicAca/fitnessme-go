package cmd

import (
	"context"
	"fitnessme/notifications/pkg/db"
	"fitnessme/notifications/pkg/repository"
	"fitnessme/notifications/pkg/services"
	"fmt"

	"github.com/robfig/cron/v3"
)

type NotificationApp struct {
	repo    repository.NotificationRepository
	service services.NotificationService
	cron    *cron.Cron
}

func New() *NotificationApp {
	handler := db.Init("localhost", "notifications", "postgres", "123", 5432)
	repo := repository.NewNotificationRepository(handler)
	service := services.NewNotificationService(repo, "smtp.gmail.com", "587", "fitnessme99@gmail.com", "cwxr tise hdqp ikzq", "fitnessme99@gmail.com")

	app := NotificationApp{
		repo:    repo,
		service: service,
		cron:    cron.New(),
	}
	return &app
}

func (n *NotificationApp) Start(ctx context.Context) error {
	_, err := n.cron.AddFunc("*/1 * * * *", func() {
		if err := n.service.ProcessNotifications(ctx); err != nil {
			fmt.Println("Failed to process notifications: ", err)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to schedule cron job: %w", err)
	}
	n.cron.Start()

	<-ctx.Done()

	n.cron.Stop()

	return nil
	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		return fmt.Errorf("context is canceled")
	// 	default:
	// 		if err := n.service.ProcessNotifications(ctx); err != nil {
	// 			fmt.Println("Failed to process notifications: ", err)
	// 		}
	// 	}
	// }
}

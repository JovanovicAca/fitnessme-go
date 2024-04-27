package repository

import (
	"fitnessme/notifications/pkg/db"
	"fitnessme/notifications/pkg/models"
)

type NotificationRepository interface {
	GetRandomNotification() (string, error)
}

type notificationRepository struct{ handler db.Handler }

func NewNotificationRepository(handler db.Handler) *notificationRepository {
	return &notificationRepository{handler: handler}
}

func (r *notificationRepository) GetRandomNotification() (string, error) {
	var notification models.Notification
	if err := r.handler.DB.Model(&models.Notification{}).Order("RANDOM()").Limit(1).Find(&notification).Error; err != nil {
		return "", err
	}
	return notification.Text, nil
}

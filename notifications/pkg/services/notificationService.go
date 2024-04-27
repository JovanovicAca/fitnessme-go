package services

import (
	"context"
	"encoding/json"
	"fitnessme/notifications/pkg/repository"
	"fmt"
	"net/http"
	"net/smtp"
)

type NotificationService interface {
	ProcessNotifications(ctx context.Context) error
	sendEmail(ctx context.Context, email, notification string) error
}

type notificationService struct {
	repo     repository.NotificationRepository
	smtpHost string
	smtpPort string
	username string
	password string
	from     string
}

func NewNotificationService(repo repository.NotificationRepository, host, port, username, password, from string) NotificationService {
	return &notificationService{
		repo:     repo,
		smtpHost: host,
		smtpPort: port,
		username: username,
		password: password,
		from:     from,
	}
}

func (n *notificationService) ProcessNotifications(ctx context.Context) error {
	emails, err := n.fetchUserEmails(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch user emails: %w", err)
	}

	notification, err := n.repo.GetRandomNotification()
	if err != nil {
		return fmt.Errorf("failed to get random notification: %w", err)
	}

	for _, email := range emails {
		if err := n.sendEmail(ctx, email, notification); err != nil {
			fmt.Printf("failed to send email to %s: %s\n", email, err)
			continue
		}
		fmt.Printf("email sent successfully to %s\n", email)
	}

	return nil
}

func (n *notificationService) sendEmail(ctx context.Context, email, notification string) error {
	addr := n.smtpHost + ":" + n.smtpPort
	auth := smtp.PlainAuth("", n.username, n.password, n.smtpHost)

	subject := "Subject: Motivation message for today\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf("<html><body>%s</body></html>", notification)
	msg := []byte("From: " + n.from + "\n" +
		"To: " + email + "\n" +
		subject + mime + body)

	err := smtp.SendMail(addr, auth, n.from, []string{email}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

func (n *notificationService) fetchUserEmails(ctx context.Context) ([]string, error) {
	userManagementServiceURL := "http://localhost:3000/user/emails"
	resp, err := http.Get(userManagementServiceURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user emails: status %d", resp.StatusCode)
	}

	var emails []string
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return nil, err
	}

	return emails, nil
}

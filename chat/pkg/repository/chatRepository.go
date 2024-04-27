package repository

import (
	"fitnessme/chat/pkg/db"
	"fitnessme/chat/pkg/models"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type ChatRepository interface {
	InsertMessage(message models.Message) error
	ChatExists(user1, user2 uuid.UUID) (string, error)
	CreateChat(chat models.Chat) error
	GetChatMessages(chatid string) ([]models.Message, error)
	UnreadMessages(user1, user2 string) (int64, error)
	ReadAll(user1, user2 string) error
	GetNewMessagesForAdmin(admin_id string) ([]models.Message, error)
}

type chatRepository struct{ handler db.Handler }

func NewChatRepository(handler db.Handler) ChatRepository {
	return &chatRepository{handler: handler}
}

func (c *chatRepository) GetNewMessagesForAdmin(admin_id string) ([]models.Message, error) {
	var messages []models.Message
	// result := c.handler.DB.Where("sent_to = ? AND status = ?", admin_id, "delivered").Order("created_at asc").Find(&messages)
	result := c.handler.DB.
		Table("messages").
		Select("DISTINCT ON (sent_by) *").
		Where("sent_to = ? AND status = ?", admin_id, "delivered").
		Order("sent_by, created_at DESC").
		Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}

func (c *chatRepository) ReadAll(user1, user2 string) error {
	err := c.handler.DB.Model(&models.Message{}).
		Where("sent_to = ? AND sent_by = ?", user1, user2).
		Update("status", "seen").Error
	return err
}

func (c *chatRepository) UnreadMessages(user1, user2 string) (int64, error) {
	var count int64
	err := c.handler.DB.Model(&models.Message{}).
		Where("sent_to = ? AND sent_by = ? AND status = ?", user1, user2, "delivered").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *chatRepository) GetChatMessages(chatid string) ([]models.Message, error) {
	var messages []models.Message
	result := c.handler.DB.Where("chat_id = ?", chatid).Order("created_at ASC").Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}

func (c *chatRepository) CreateChat(chat models.Chat) error {
	if err := c.handler.DB.Create(&chat).Error; err != nil {
		return errors.Wrap(err, "failed to save chat")
	}
	return nil
}

func (c *chatRepository) ChatExists(user1, user2 uuid.UUID) (string, error) {
	var chatID string
	if err := c.handler.DB.Model(&models.Chat{}).
		Where("(user_1 = ? AND user_2 = ?) OR (user_1 = ? AND user_2 = ?)", user1, user2, user2, user1).
		Pluck("id", &chatID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return chatID, nil
}

func (c *chatRepository) InsertMessage(message models.Message) error {
	if err := c.handler.DB.Create(&message).Error; err != nil {
		return errors.Wrap(err, "failed to insert message")
	}
	return nil
}

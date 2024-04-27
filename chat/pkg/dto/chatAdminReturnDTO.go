package dto

import "github.com/google/uuid"

type ChatAdminReturnDTO struct {
	UserName string    `json:"name"`
	SentBy   uuid.UUID `json:"sent_by"`
}

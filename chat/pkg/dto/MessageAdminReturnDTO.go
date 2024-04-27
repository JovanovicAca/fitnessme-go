package dto

import (
	"github.com/google/uuid"
)

type MessageAdminReturnDTO struct {
	Id        uuid.UUID `json:"id"`
	Text      string    `json:"text"`
	User      string    `json:"user"`
	Status    string    `json:"status"`
	SentBy    uuid.UUID `json:"sent_by"`
	SentTo    uuid.UUID `json:"sent_to"`
	RepliedOn string    `json:"replied_on" `
	ChatId    string    `json:"chat_id"`
}

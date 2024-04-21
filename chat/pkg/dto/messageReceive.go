package dto

import (
	"github.com/google/uuid"
)

type MessageReceive struct {
	Text      string    `json:"text"`
	SentBy    uuid.UUID `json:"sent_by"`
	RepliedOn string    `json:"replied_on"`
	ChatId    uuid.UUID `json:"chat_id"`
}

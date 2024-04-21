package dto

import (
	"github.com/google/uuid"
)

type ChatDTO struct {
	User1 uuid.UUID `json:"user_1"`
	User2 uuid.UUID `json:"user_2"`
}

package domain

import (
	"time"

	"github.com/google/uuid"
)

// Use encoding/json/v2 instead of encoding/json to omit zero values
// 使いたいだけです

type Event struct {
	ID      uuid.UUID `json:"id"`
	Title   string    `json:"title"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at,omitzero"`
}

type Comment struct {
	ID        uuid.UUID `json:"id"`
	UserID    string    `json:"user_id"`
	EventID   uuid.UUID `json:"event_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

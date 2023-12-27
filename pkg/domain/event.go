package domain

import (
	"errors"
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

func (e *Event) Validate() error {
	if e.Title == "" {
		return errors.New("title is required")
	}
	if e.StartAt.IsZero() {
		return errors.New("start_at is required")
	}
	return nil
}

package domain

import (
	"time"

	"github.com/google/uuid"
)

type StreamObject interface {
	Type() StreamObjectType
}

type StreamObjectType int

const (
	StreamObjectTypeComment StreamObjectType = iota
)

func (s StreamObjectType) String() string {
	switch s {
	case StreamObjectTypeComment:
		return "comment"
	default:
		return "unknown"
	}
}

type Comment struct {
	ID        uuid.UUID `json:"id"`
	UserID    string    `json:"user_id"`
	EventID   uuid.UUID `json:"event_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Comment) Type() StreamObjectType {
	return StreamObjectTypeComment
}

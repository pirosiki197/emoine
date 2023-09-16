package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Event struct {
	bun.BaseModel `bun:"events"`

	ID      uuid.UUID `bun:"id,pk,type:uuid"`
	Title   string    `bun:"title,notnull,nullzero"`
	StartAt time.Time `bun:"start_at,notnull,nullzero,default:current_timestamp"`
	EndAt   time.Time `bun:"end_at,nullzero"`
}

type Comment struct {
	bun.BaseModel `bun:"comments"`

	ID        uuid.UUID `bun:"id,pk,type:uuid"`
	UserID    string    `bun:"user_id,notnull,nullzero"`
	EventID   uuid.UUID `bun:"event_id,notnull,nullzero"`
	Text      string    `bun:"text,notnull,nullzero"`
	CreatedAt time.Time `bun:"created_at,notnull,nullzero,default:current_timestamp"`
}

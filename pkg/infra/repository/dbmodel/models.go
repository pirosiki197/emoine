package dbmodel

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Event struct {
	bun.BaseModel `bun:"events"`

	ID      uuid.UUID    `bun:",pk,type:varchar(36)"`
	Title   string       `bun:",notnull,nullzero"`
	StartAt time.Time    `bun:",notnull,nullzero,default:current_timestamp"`
	EndAt   sql.NullTime `bun:",nullzero"`
}

type Comment struct {
	bun.BaseModel `bun:"comments"`

	ID        uuid.UUID `bun:",pk,type:varchar(36)"`
	UserID    string    `bun:",notnull,nullzero"`
	EventID   uuid.UUID `bun:",notnull,nullzero,type:varchar(36)"`
	Text      string    `bun:",notnull,nullzero"`
	CreatedAt time.Time `bun:",notnull,nullzero,default:current_timestamp"`
}

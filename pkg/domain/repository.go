package domain

import (
	"context"
)

type Repository interface {
	CreateEvent(ctx context.Context, event *Event) error
	GetEvents(ctx context.Context) ([]Event, error)
	GetEvent(ctx context.Context, id string) (*Event, error)

	CreateComment(ctx context.Context, comment *Comment) error
	GetEventComments(ctx context.Context, eventID string) ([]Comment, error)
}

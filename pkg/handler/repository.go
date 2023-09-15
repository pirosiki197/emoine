package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/pirosiki197/emoine/pkg/model"
)

type Repository interface {
	CreateEvent(ctx context.Context, event *model.Event) error
	GetEvents(ctx context.Context) ([]model.Event, error)
	GetEvent(ctx context.Context, id uuid.UUID) (*model.Event, error)

	SendComment(ctx context.Context, comment *model.Comment) error
	GetEventComments(ctx context.Context, eventID string) ([]model.Comment, error)
}

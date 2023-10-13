package handler

import (
	"context"

	"github.com/pirosiki197/emoine/pkg/domain"
)

type Repository interface {
	CreateEvent(ctx context.Context, event *domain.Event) error
	GetEvents(ctx context.Context) ([]domain.Event, error)
	GetEvent(ctx context.Context, id string) (*domain.Event, error)

	CreateComment(ctx context.Context, comment *domain.Comment) error
	GetEventComments(ctx context.Context, eventID string) ([]domain.Comment, error)
}

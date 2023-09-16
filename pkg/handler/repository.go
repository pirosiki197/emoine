package handler

import (
	"context"

	"github.com/pirosiki197/emoine/pkg/model"
)

type Repository interface {
	CreateEvent(ctx context.Context, event *model.Event) error
	GetEvents(ctx context.Context) ([]model.Event, error)
	GetEvent(ctx context.Context, id string) (*model.Event, error)

	CreateComment(ctx context.Context, comment *model.Comment) error
	GetEventComments(ctx context.Context, eventID string) ([]model.Comment, error)
}

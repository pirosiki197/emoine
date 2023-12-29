package dbmodel

import (
	"database/sql"

	"github.com/pirosiki197/emoine/pkg/domain"
	"github.com/samber/lo"
)

func FromDomainEvent(e *domain.Event) *Event {
	if e == nil {
		return nil
	}
	return &Event{
		ID:      e.ID,
		Title:   e.Title,
		StartAt: e.StartAt,
		EndAt:   lo.Ternary(e.EndAt.IsZero(), sql.NullTime{}, sql.NullTime{Time: e.EndAt, Valid: true}),
	}
}

func (e *Event) ToDomain() *domain.Event {
	if e == nil {
		return nil
	}
	return &domain.Event{
		ID:      e.ID,
		Title:   e.Title,
		StartAt: e.StartAt,
		EndAt:   e.EndAt.Time,
	}
}

func FromDomainComment(c *domain.Comment) *Comment {
	if c == nil {
		return nil
	}
	return &Comment{
		ID:        c.ID,
		UserID:    c.UserID,
		EventID:   c.EventID,
		Text:      c.Text,
		CreatedAt: c.CreatedAt,
	}
}

func (c *Comment) ToDomain() *domain.Comment {
	if c == nil {
		return nil
	}
	return &domain.Comment{
		ID:        c.ID,
		UserID:    c.UserID,
		EventID:   c.EventID,
		Text:      c.Text,
		CreatedAt: c.CreatedAt,
	}
}

package dbmodel

import "github.com/pirosiki197/emoine/pkg/domain"

func FromDomainEvent(e *domain.Event) *Event {
	if e == nil {
		return nil
	}
	return &Event{
		ID:      e.ID,
		Title:   e.Title,
		StartAt: e.StartAt,
		EndAt:   e.EndAt,
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
		EndAt:   e.EndAt,
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

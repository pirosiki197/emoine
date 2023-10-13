package repository

import (
	"context"

	"github.com/pirosiki197/emoine/pkg/domain"
	"github.com/pirosiki197/emoine/pkg/repository/dbmodel"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
)

func (r *Repository) CreateEvent(ctx context.Context, event *domain.Event) error {
	e := dbmodel.FromDomainEvent(event)

	_, err := r.db.
		NewInsert().
		Model(e).
		Exec(ctx)

	return err
}

func (r *Repository) GetEvents(ctx context.Context) ([]domain.Event, error) {
	var events []dbmodel.Event

	err := r.db.
		NewSelect().
		Model(&events).
		OrderExpr("? DESC", bun.Ident("start_at")).
		Scan(ctx)

	return lo.Map(events, func(e dbmodel.Event, _ int) domain.Event {
		return *e.ToDomain()
	}), err
}

func (r *Repository) GetEvent(ctx context.Context, id string) (*domain.Event, error) {
	var event dbmodel.Event

	err := r.db.
		NewSelect().
		Model(&event).
		Where("? = ?", bun.Ident("id"), id).
		Scan(ctx)

	return event.ToDomain(), err
}

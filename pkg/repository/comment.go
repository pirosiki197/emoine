package repository

import (
	"context"

	"github.com/pirosiki197/emoine/pkg/domain"
	"github.com/pirosiki197/emoine/pkg/repository/dbmodel"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
)

func (r *Repository) CreateComment(ctx context.Context, comment *domain.Comment) error {
	c := dbmodel.FromDomainComment(comment)

	_, err := r.db.
		NewInsert().
		Model(c).
		Exec(ctx)

	return err
}

func (r *Repository) GetEventComments(ctx context.Context, eventID string) ([]domain.Comment, error) {
	var comments []dbmodel.Comment

	err := r.db.
		NewSelect().
		Model(&comments).
		Where("? = ?", bun.Ident("event_id"), eventID).
		OrderExpr("? ASC", bun.Ident("created_at")).
		Scan(ctx)

	return lo.Map(comments, func(c dbmodel.Comment, _ int) domain.Comment {
		return *c.ToDomain()
	}), err
}

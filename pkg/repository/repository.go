package repository

import (
	"context"

	"github.com/pirosiki197/emoine/pkg/model"
	"github.com/uptrace/bun"
)

// Repository is a struct that implements the handler.Repository interface.
type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	setUpDB(db)
	return &Repository{
		db: db,
	}
}

func setUpDB(db *bun.DB) {
	ctx := context.Background()
	_, err := db.NewCreateTable().Model((*model.Event)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		panic(err)
	}
	_, err = db.NewCreateTable().Model((*model.Comment)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		panic(err)
	}
}

func (r *Repository) CreateEvent(ctx context.Context, event *model.Event) error {
	_, err := r.db.NewInsert().Model(event).Exec(ctx)

	return err
}

func (r *Repository) GetEvents(ctx context.Context) ([]model.Event, error) {
	var events []model.Event

	err := r.db.NewSelect().Model(&events).Order("start_at DESC").Scan(ctx)

	return events, err
}

func (r *Repository) GetEvent(ctx context.Context, id string) (*model.Event, error) {
	var event model.Event

	err := r.db.NewSelect().Model(&event).Where("id = ?", id).Scan(ctx)

	return &event, err
}

func (r *Repository) CreateComment(ctx context.Context, comment *model.Comment) error {
	_, err := r.db.NewInsert().Model(comment).Exec(ctx)

	return err
}

func (r *Repository) GetEventComments(ctx context.Context, eventID string) ([]model.Comment, error) {
	var comments []model.Comment

	err := r.db.NewSelect().Model(&comments).Where("event_id = ?", eventID).Scan(ctx)

	return comments, err
}

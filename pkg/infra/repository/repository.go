package repository

import (
	"context"

	"github.com/pirosiki197/emoine/pkg/infra/repository/dbmodel"
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
	_, err := db.NewCreateTable().Model((*dbmodel.Event)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		panic(err)
	}
	_, err = db.NewCreateTable().Model((*dbmodel.Comment)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		panic(err)
	}
}

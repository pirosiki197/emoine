package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/pirosiki197/emoine/pkg/repository/dbmodel"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

// Repository is a struct that implements the handler.Repository interface.
type Repository struct {
	db *bun.DB
}

func NewRepository() *Repository {
	db := newDB()
	setUpDB(db)
	return &Repository{
		db: db,
	}
}

func newDB() *bun.DB {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}

	conf := mysql.Config{
		User:      "root",
		Passwd:    "example",
		Net:       "tcp",
		Addr:      "db:3306",
		DBName:    "emoine",
		Loc:       jst,
		ParseTime: true,
	}
	var sqldb *sql.DB
	for i := 0; i < 10; i++ {
		time.Sleep(2 * time.Second)
		sqldb, _ = sql.Open("mysql", conf.FormatDSN())
		if err := sqldb.Ping(); err == nil {
			break
		}
		if i == 9 {
			panic(err)
		}
	}
	db := bun.NewDB(sqldb, mysqldialect.New())

	return db
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

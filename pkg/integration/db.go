package integration

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/pirosiki197/emoine/pkg/infra/repository/dbmodel"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

type DB struct {
	db *bun.DB
}

func NewDB(t *testing.T) *DB {
	t.Helper()
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	conf := mysql.Config{
		User:                 "root",
		Passwd:               "example",
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               "emoine",
		Loc:                  jst,
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	db, err := connectDB(&conf)
	if err != nil {
		t.Fatal(err)
	}
	return &DB{
		db: db,
	}
}

func connectDB(conf *mysql.Config) (db *bun.DB, err error) {
	sqldb, err := sql.Open("mysql", conf.FormatDSN())
	if err != nil {
		return nil, err
	}
	if err := sqldb.Ping(); err != nil {
		return nil, err
	}
	return bun.NewDB(sqldb, mysqldialect.New()), nil
}

func (d *DB) Cleanup() {
	ctx := context.Background()
	_, err := d.db.NewDropTable().Model((*dbmodel.Event)(nil)).IfExists().Exec(ctx)
	if err != nil {
		panic(err)
	}
	_, err = d.db.NewDropTable().Model((*dbmodel.Comment)(nil)).IfExists().Exec(ctx)
	if err != nil {
		panic(err)
	}
}

func (d *DB) BunDB() *bun.DB {
	return d.db
}

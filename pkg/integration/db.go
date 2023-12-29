package integration

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/pirosiki197/emoine/pkg/infra/repository/dbmodel"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

type DB struct {
	db  *bun.DB
	rdb *redis.Client
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
	rdb := connectRedis(t)
	return &DB{
		db:  db,
		rdb: rdb,
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

func connectRedis(t *testing.T) *redis.Client {
	t.Helper()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		t.Fatal(err)
	}
	return rdb
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

func (d *DB) Redis() *redis.Client {
	return d.rdb
}

package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/pirosiki197/emoine/pkg/handler"
	"github.com/pirosiki197/emoine/pkg/infra/proto/api/v1/apiv1connect"
	"github.com/pirosiki197/emoine/pkg/infra/repository"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(l)

	repo := repository.NewRepository(ConnectBunDB())
	rdb := ConnectRedis()
	handler := handler.NewHandler(repo, rdb)

	mux := http.NewServeMux()
	mux.Handle(apiv1connect.NewAPIServiceHandler(handler))
	c := cors.AllowAll()

	slog.Info("server started on :8080")
	log.Fatal(http.ListenAndServe(
		":8080",
		h2c.NewHandler(c.Handler(mux), &http2.Server{}),
	))
}

func ConnectBunDB() *bun.DB {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	conf := mysql.Config{
		User:                 "root",
		Passwd:               "password",
		Net:                  "tcp",
		Addr:                 "db:3306",
		DBName:               "emoine",
		Loc:                  jst,
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	for i := 0; i < 10; i++ {
		db, err := connectDB(&conf)
		if err != nil {
			slog.Warn("failed to connect to db", slog.String("err", err.Error()), slog.Int("retry", i))
			time.Sleep(5 * time.Second)
			continue
		}
		return db
	}
	panic(err)
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

func ConnectRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
}

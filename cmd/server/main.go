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
	"github.com/rs/cors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(logger)

	repo := repository.NewRepository(ConnectDB())
	handler := handler.NewHandlre(repo)

	mux := http.NewServeMux()
	mux.Handle(apiv1connect.NewAPIServiceHandler(handler))
	c := cors.AllowAll()

	slog.Info("server started on :8080")
	log.Fatal(http.ListenAndServe(
		":8080",
		h2c.NewHandler(c.Handler(mux), &http2.Server{}),
	))
}

func ConnectDB() *bun.DB {
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
	sqldb, err := sql.Open("mysql", conf.FormatDSN())
	if err != nil {
		panic(err)
	}
	return bun.NewDB(sqldb, mysqldialect.New())
}

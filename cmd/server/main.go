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
	"github.com/pirosiki197/emoine/pkg/proto/api/v1/apiv1connect"
	"github.com/pirosiki197/emoine/pkg/repository"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	repo := repository.NewRepository(newDB())
	handler := handler.NewHandlre(repo)
	mux := http.NewServeMux()
	mux.Handle(apiv1connect.NewAPIServiceHandler(handler))

	logger.Info("server started on :8080")
	log.Fatal(http.ListenAndServe(
		":8080",
		h2c.NewHandler(mux, &http2.Server{}),
	))
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

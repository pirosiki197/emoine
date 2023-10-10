package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/pirosiki197/emoine/pkg/handler"
	"github.com/pirosiki197/emoine/pkg/proto/api/v1/apiv1connect"
	"github.com/pirosiki197/emoine/pkg/repository"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(logger)

	repo := repository.NewRepository()
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

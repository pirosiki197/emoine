package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/pirosiki197/emoine/pkg/gen/api/v1/apiv1connect"
	"github.com/pirosiki197/emoine/pkg/handler"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	handler := handler.NewHandlre(nil)
	mux := http.NewServeMux()
	mux.Handle(apiv1connect.NewAPIServiceHandler(handler))
	logger.Info("server started on :8080")
	log.Fatal(http.ListenAndServe(
		":8080",
		h2c.NewHandler(mux, &http2.Server{}),
	))
}

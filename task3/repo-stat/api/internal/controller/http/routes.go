package http

import (
	"log/slog"
	"net/http"
	"repo-stat/api/internal/adapter/processor"
	"repo-stat/api/internal/usecase"
)

func AddRoutes(mux *http.ServeMux, log *slog.Logger, ping *usecase.Ping, processorClient *processor.Client) {
	mux.Handle("GET /api/ping", NewPingHandler(log, ping))
	mux.Handle("GET /api/repositories/info", NewRepoInfoHandler(log, processorClient))
}

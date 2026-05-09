package http

import (
	"log/slog"
	"net/http"
	"repo-stat/api/internal/usecase"
)

func AddRoutes(mux *http.ServeMux, log *slog.Logger, ping *usecase.Ping, repoInfo *usecase.RepoInfo, subscriptionManager *usecase.SubscriptionManager) {
	mux.Handle("GET /api/ping", NewPingHandler(log, ping))
	mux.Handle("GET /api/repositories/info", NewRepoInfoHandler(log, repoInfo))
	mux.Handle("POST /api/subscriptions", NewAddSubscriptionHandler(log, subscriptionManager))
	mux.Handle("GET /api/subscriptions", NewListSubscriptionsHandler(log, subscriptionManager))
	mux.Handle("DELETE /api/subscriptions/{owner}/{repo}", NewDeleteSubscriptionHandler(log, subscriptionManager))
	mux.Handle("GET /api/subscriptions/info", NewSubscriptionsInfoHandler(log, repoInfo))
}

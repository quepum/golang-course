package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
	"strings"
)

// NewAddSubscriptionHandler
// @Summary Add a new repository subscription
// @Description Adds a subscription to a GitHub repository. Validates that the repository exists on GitHub before saving to the database.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body dto.AddSubscriptionRequest true "Owner and Repository name"
// @Success 201 {object} dto.AddSubscriptionResponse "Subscription added successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid JSON format or repository not found on GitHub"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /subscriptions [post]
func NewAddSubscriptionHandler(log *slog.Logger, subUC *usecase.SubscriptionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req dto.AddSubscriptionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Warn("invalid json in add subscription request", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "invalid json"})
			return
		}

		log.Info("adding subscription", "owner", req.Owner, "repo", req.Repo)

		err := subUC.AddSubscription(r.Context(), req.Owner, req.Repo)
		if err != nil {
			log.Error("failed to add subscription via usecase", "owner", req.Owner, "repo", req.Repo, "error", err)

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "failed to add subscription"})
			return
		}

		log.Info("subscription added successfully", "owner", req.Owner, "repo", req.Repo)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(dto.AddSubscriptionResponse{Status: true, Message: "added"})
	}
}

// NewListSubscriptionsHandler
// @Summary Get list of all active subscriptions
// @Description Returns a list of all repositories currently subscribed by the user from the database.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 200 {object} dto.ListSubscriptionsResponse "List of subscribed repositories"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /subscriptions [get]
func NewListSubscriptionsHandler(log *slog.Logger, subUC *usecase.SubscriptionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		log.Debug("listing subscriptions")

		subs, err := subUC.ListSubscriptions(r.Context())
		if err != nil {
			log.Error("failed to list subscriptions via usecase", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "internal error"})
			return
		}

		var items []dto.SubscriptionItem
		for _, s := range subs {
			items = append(items, dto.SubscriptionItem{
				Owner: s.Owner,
				Repo:  s.Repo,
			})
		}

		log.Info("subscriptions listed", "count", len(items))
		json.NewEncoder(w).Encode(dto.ListSubscriptionsResponse{Subscriptions: items})
	}
}

// NewDeleteSubscriptionHandler
// @Summary Remove a specific subscription
// @Description Removes a subscription for a given owner and repository name from the database.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param owner path string true "GitHub Owner username"
// @Param repo path string true "Repository name"
// @Success 200 {object} dto.RemoveSubscriptionResponse "Subscription removed successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid parameters"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /subscriptions/{owner}/{repo} [delete]
func NewDeleteSubscriptionHandler(log *slog.Logger, subUC *usecase.SubscriptionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		path := strings.TrimPrefix(r.URL.Path, "/api/subscriptions/")
		parts := strings.Split(path, "/")

		if len(parts) != 2 {
			log.Warn("invalid delete subscription url format", "path", r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "invalid url format"})
			return
		}

		owner, repo := parts[0], parts[1]
		log.Info("removing subscription", "owner", owner, "repo", repo)

		err := subUC.RemoveSubscription(r.Context(), owner, repo)
		if err != nil {
			log.Error("failed to remove subscription via usecase", "owner", owner, "repo", repo, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "failed to remove subscription"})
			return
		}

		log.Info("subscription removed successfully", "owner", owner, "repo", repo)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(dto.RemoveSubscriptionResponse{Status: true})
	}
}

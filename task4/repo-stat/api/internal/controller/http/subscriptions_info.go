package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

// NewSubscriptionsInfoHandler
// @Summary Get statistics for all subscribed repositories
// @Description Fetches real-time statistics (stars, forks, description) from GitHub for every repository in the subscription list.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 200 {object} dto.SubscriptionsInfoResponse "Statistics for all subscribed repos"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /subscriptions/info [get]
func NewSubscriptionsInfoHandler(log *slog.Logger, repoInfoUC *usecase.RepoInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		domainRepos, err := repoInfoUC.GetSubscribedStats(r.Context())
		if err != nil {
			log.Error("failed to get subscribed stats", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "internal server error"})
			return
		}

		var stats []dto.RepoStatItem
		for _, repo := range domainRepos {
			stats = append(stats, dto.RepoStatItem{
				FullName:    repo.FullName,
				Description: repo.Description,
				Stars:       repo.Stars,
				Forks:       repo.Forks,
				CreatedAt:   repo.CreatedAt,
			})
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(dto.SubscriptionsInfoResponse{Stats: stats})
	}
}

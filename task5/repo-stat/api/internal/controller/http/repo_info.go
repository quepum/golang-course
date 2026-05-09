package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/usecase"
	"strings"

	"repo-stat/api/internal/dto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewRepoInfoHandler
// @Summary      Get repository information
// @Description  Returns statistics for a given GitHub repository URL
// @Tags         repositories
// @Produce      json
// @Param        url   query     string  true  "GitHub Repository URL"
// @Success      200   {object}  dto.RepoInfoResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      404   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /repositories/info [get]
func NewRepoInfoHandler(log *slog.Logger, repoInfoUC *usecase.RepoInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		urlStr := r.URL.Query().Get("url")
		if urlStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ErrorResponse{Error: `missing required query parameter "url"`})
			return
		}

		if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "invalid url format"})
			return
		}

		domainRepo, err := repoInfoUC.Execute(r.Context(), urlStr)
		if err != nil {
			handleGrpcError(w, log, err)
			return
		}

		response := dto.RepoInfoResponse{
			FullName:    domainRepo.FullName,
			Description: domainRepo.Description,
			Stars:       domainRepo.Stars,
			Forks:       domainRepo.Forks,
			CreatedAt:   domainRepo.CreatedAt,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func handleGrpcError(w http.ResponseWriter, log *slog.Logger, err error) {
	st, ok := status.FromError(err)
	if ok {
		switch st.Code() {
		case codes.NotFound:
			w.WriteHeader(http.StatusNotFound)
			resp := dto.ErrorResponse{Error: "repository not found"}
			_ = json.NewEncoder(w).Encode(resp)
			return

		case codes.ResourceExhausted:
			w.WriteHeader(http.StatusTooManyRequests)
			resp := dto.ErrorResponse{Error: "github rate limit exceeded"}
			_ = json.NewEncoder(w).Encode(resp)
			return

		case codes.InvalidArgument:
			w.WriteHeader(http.StatusBadRequest)
			resp := dto.ErrorResponse{Error: "invalid argument: " + st.Message()}
			_ = json.NewEncoder(w).Encode(resp)
			return

		default:
			log.Error("gRPC error", "code", st.Code(), "message", st.Message())
			w.WriteHeader(http.StatusInternalServerError)
			resp := dto.ErrorResponse{Error: "internal server error"}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
	}

	log.Error("unknown error", "error", err)
	w.WriteHeader(http.StatusInternalServerError)
	resp := dto.ErrorResponse{Error: "internal server error"}
	_ = json.NewEncoder(w).Encode(resp)
}

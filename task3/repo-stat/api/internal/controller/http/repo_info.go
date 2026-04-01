package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"repo-stat/api/internal/adapter/processor"
	"repo-stat/api/internal/dto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewRepoInfoHandler(log *slog.Logger, procClient *processor.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		urlStr := r.URL.Query().Get("url")
		if urlStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp := dto.ErrorResponse{Error: `missing required query parameter "url"`}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
			w.WriteHeader(http.StatusBadRequest)
			resp := dto.ErrorResponse{Error: "invalid url format: must start with http:// or https://"}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		resp, err := procClient.GetRepoInfo(r.Context(), urlStr)
		if err != nil {
			handleGrpcError(w, log, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("failed to encode response", "error", err)
		}
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

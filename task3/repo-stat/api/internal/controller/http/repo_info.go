package http

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"repo-stat/api/internal/adapter/processor"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewRepoInfoHandler(log *slog.Logger, procClient *processor.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlStr := r.URL.Query().Get("url")
		if urlStr == "" {
			http.Error(w, `missing required query parameter "url"`, http.StatusBadRequest)
			return
		}

		if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
			http.Error(w, "invalid url format: must start with http:// or https://", http.StatusBadRequest)
			return
		}

		resp, err := procClient.GetRepoInfo(r.Context(), urlStr)
		if err != nil {
			handleGrpcError(w, log, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
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
			http.Error(w, "repository not found", http.StatusNotFound)
			return
		case codes.ResourceExhausted:
			http.Error(w, "github rate limit exceeded", http.StatusTooManyRequests)
			return
		case codes.InvalidArgument:
			http.Error(w, fmt.Sprintf("invalid argument: %s", st.Message()), http.StatusBadRequest)
			return
		default:
			log.Error("gRPC error", "code", st.Code(), "message", st.Message())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	log.Error("unknown error", "error", err)
	http.Error(w, "internal server error", http.StatusInternalServerError)
}

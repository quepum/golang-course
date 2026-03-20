package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"task2/domain"
)

type HttpServer struct {
	uc domain.UseCase
}

func NewHTTPServer(uc domain.UseCase) *HttpServer {
	return &HttpServer{
		uc: uc,
	}
}

// GetRepoInfo retrieves information about a GitHub repository.
// @Summary      Get repository information
// @Description  Retrieves detailed information about a GitHub repository including name, description, stars, forks, and creation date.
// @Tags         repos
// @Accept       json
// @Produce      json
// @Param        owner  path      string  true  "Repository Owner username (e.g. golang)"
// @Param        repo   path      string  true  "Repository Name (e.g. go)"
// @Success      200    {object}  domain.RepoInfo  "Successful response with repository details"
// @Failure      400    {string}  string           "Invalid input parameters: owner or repo is empty"
// @Failure      404    {string}  string           "Repository not found on GitHub"
// @Failure      500    {string}  string           "Internal server error: failed to fetch data from collector"
// @Router       /repos/{owner}/{repo} [get]
func (h *HttpServer) GetRepoInfo(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) != 3 || parts[0] != "repos" {
		http.Error(w, "invalid format", http.StatusBadRequest)
		return
	}

	owner, repo := parts[1], parts[2]
	info, err := h.uc.GetRepoInfo(r.Context(), owner, repo)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRepoNotFound):
			http.Error(w, "repository not found", http.StatusNotFound)
			return

		case errors.Is(err, domain.ErrInvalidInput):
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

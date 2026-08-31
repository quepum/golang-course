package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"task2/collector/domain"
	"time"
)

type githubResponse struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	StargazersCount int    `json:"stargazers_count"` // Имя поля как в API
	ForksCount      int    `json:"forks_count"`      // Имя поля как в API
	CreatedAt       string `json:"created_at"`       // Приходит строкой
}

type githubRepo struct {
	client      *http.Client
	standardUrl string
}

func NewGitHubRepo() domain.Repository {
	return &githubRepo{
		client:      &http.Client{Timeout: 10 * time.Second},
		standardUrl: "https://api.github.com/repos",
	}
}

func (r *githubRepo) GetRepoInfo(ctx context.Context, owner, repo string) (*domain.RepoInfo, error) {
	url := fmt.Sprintf("%s/%s/%s", r.standardUrl, owner, repo)

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	response, err := r.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return nil, domain.ErrRepoNotFound
		}
		return nil, fmt.Errorf("github api error: status %d, body: %s", response.StatusCode, string(body))
	}

	var rawDTO githubResponse
	if err := json.Unmarshal(body, &rawDTO); err != nil {
		return nil, fmt.Errorf("error decoding response to adapter model: %w", err)
	}

	return &domain.RepoInfo{
		Name:        rawDTO.Name,
		Description: rawDTO.Description,
		Stars:       rawDTO.StargazersCount,
		Forks:       rawDTO.ForksCount,
		CreatedAt:   rawDTO.CreatedAt,
	}, nil
}

package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"task2/domain"
	"time"
)

type githubRepo struct {
	client      *http.Client
	standardUrl string
}

func NewGitHubRepo() domain.Repository {
	newRepo := &githubRepo{
		client:      &http.Client{Timeout: 10 * time.Second},
		standardUrl: "https://api.github.com/repos",
	}

	return newRepo
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

	if response.StatusCode != http.StatusOK {
		_, _ = io.ReadAll(response.Body)
		if response.StatusCode == http.StatusNotFound {
			return nil, domain.ErrRepoNotFound
		}

		return nil, fmt.Errorf("github api error: status %d", response.StatusCode)
	}

	var result domain.RepoInfo
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil

}

package usecase

import (
	"context"
	"fmt"
	"repo-stat/processor/internal/adapter/collector"
	collectorv1 "repo-stat/proto/collector"
	"strings"
)

type RepoInfo struct {
	collClient *collector.Client
}

func NewRepoInfo(collClient *collector.Client) *RepoInfo {
	return &RepoInfo{collClient: collClient}
}

func (u *RepoInfo) Execute(ctx context.Context, url string) (*collectorv1.RepoInfoResponse, error) {
	owner, repo, err := parseURL(url)
	if err != nil {
		return nil, err
	}

	return u.collClient.GetRepoInfo(ctx, owner, repo)
}

func parseURL(url string) (owner, repo string, err error) {
	url = strings.TrimSuffix(url, "/")

	if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "https") {
		return "", "", fmt.Errorf("invalid url: %s", url)
	}

	parts := strings.Split(url, "/")

	if len(parts) < 5 {
		return "", "", fmt.Errorf("invalid url: %s", url)
	}

	if parts[2] != "github.com" {
		return "", "", fmt.Errorf("invalid url: %s", url)
	}

	owner = parts[3]
	repo = parts[4]

	if owner == "" || repo == "" {
		return "", "", fmt.Errorf("empty owner or repo: %s", url)
	}

	return owner, repo, nil
}

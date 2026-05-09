package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type RepoInfo struct {
	procClient ProcessorClient
}

func NewRepoInfo(procClient ProcessorClient) *RepoInfo {
	return &RepoInfo{procClient: procClient}
}

func (r *RepoInfo) Execute(ctx context.Context, url string) (*domain.Repository, error) {
	return r.procClient.GetRepoInfo(ctx, url)
}

func (r *RepoInfo) GetSubscribedStats(ctx context.Context) ([]domain.Repository, error) {
	return r.procClient.GetSubscribedStats(ctx)
}

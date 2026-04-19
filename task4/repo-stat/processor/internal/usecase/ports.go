package usecase

import (
	"context"
	"repo-stat/processor/internal/domain"
)

type CollectorClient interface {
	GetRepoInfo(ctx context.Context, owner, repo string) (*domain.Repository, error)
	GetSubscribedStats(ctx context.Context) ([]*domain.Repository, error)
}

package usecase

import (
	"context"
	"repo-stat/collector/internal/domain"
)

type GitHubClient interface {
	GetRepoInfo(ctx context.Context, owner, repo string) (*domain.Repository, error)
}

type SubscriptionLister interface {
	GetActiveSubscriptions(ctx context.Context) ([]*domain.Subscription, error)
}

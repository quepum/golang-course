package usecase

import (
	"context"
	"repo-stat/subscriber/internal/domain"
)

type SubscriptionRepository interface {
	AddSubscription(ctx context.Context, sub *domain.Subscription) error
	ListSubscriptions(ctx context.Context) ([]domain.Subscription, error)
	RemoveSubscription(ctx context.Context, owner, repo string) error
}

type GitHubChecker interface {
	CheckRepoExists(ctx context.Context, owner, repo string) error
}

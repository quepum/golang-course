package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type Pinger interface {
	Ping(ctx context.Context) domain.PingStatus
	GetName() string
}

type ProcessorClient interface {
	GetRepoInfo(ctx context.Context, url string) (*domain.Repository, error)
	GetSubscribedStats(ctx context.Context) ([]domain.Repository, error)
}

type SubscriberClient interface {
	AddSubscription(ctx context.Context, owner, repo string) error
	RemoveSubscription(ctx context.Context, owner, repo string) error
	ListSubscriptions(ctx context.Context) ([]domain.Subscription, error)
}

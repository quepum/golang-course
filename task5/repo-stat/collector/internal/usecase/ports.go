package usecase

import (
	"context"
	"repo-stat/collector/internal/domain"
	platformKafka "repo-stat/platform/kafka"
)

type GitHubClient interface {
	GetRepoInfo(ctx context.Context, owner, repo string) (*domain.Repository, error)
}

type SubscriptionLister interface {
	GetActiveSubscriptions(ctx context.Context) ([]*domain.Subscription, error)
}

type KafkaProducer interface {
	SendFetchResult(ctx context.Context, res platformKafka.FetchResult) error
	SendFetchRequest(ctx context.Context, req platformKafka.FetchRequest) error
}

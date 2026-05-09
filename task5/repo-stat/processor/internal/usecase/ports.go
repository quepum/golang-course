package usecase

import (
	"context"
	"repo-stat/platform/kafka"
	"repo-stat/processor/internal/domain"
)

type CollectorClient interface {
	GetRepoInfo(ctx context.Context, owner, repo string) (*domain.Repository, error)
	GetSubscribedStats(ctx context.Context) ([]*domain.Repository, error)
}

type KafkaProducer interface {
	SendFetchRequest(ctx context.Context, req kafka.FetchRequest) error
}

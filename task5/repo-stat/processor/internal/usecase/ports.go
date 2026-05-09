package usecase

import (
	"context"
	"repo-stat/platform/kafka"
	"repo-stat/processor/internal/adapter/storage"
)

type KafkaProducer interface {
	SendFetchRequest(ctx context.Context, req kafka.FetchRequest) error
}

type RepoStorage interface {
	GetInfoByOwnerAndRepo(ctx context.Context, arg storage.GetInfoByOwnerAndRepoParams) (storage.Repository, error)
	ListAllStats(ctx context.Context) ([]storage.Repository, error)
}

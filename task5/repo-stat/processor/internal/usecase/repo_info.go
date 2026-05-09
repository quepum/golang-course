package usecase

import (
	"context"
	"errors"
	"fmt"
	"repo-stat/platform/kafka"
	"repo-stat/processor/internal/adapter/storage"
	"repo-stat/processor/internal/domain"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type RepoInfo struct {
	storage       *storage.Repo
	kafkaProducer KafkaProducer
}

func NewRepoInfo(storage *storage.Repo, kafkaProducer KafkaProducer) *RepoInfo {
	return &RepoInfo{
		storage:       storage,
		kafkaProducer: kafkaProducer,
	}
}

func (u *RepoInfo) Execute(ctx context.Context, url string) (*domain.Repository, error) {
	owner, repo, err := parseURL(url)
	if err != nil {
		return nil, err
	}

	db, err := u.storage.GetInfoByOwnerAndRepo(ctx, storage.GetInfoByOwnerAndRepoParams{
		Owner: owner,
		Repo:  repo,
	})

	if err == nil {
		return &domain.Repository{
			FullName:    db.FullName.String,
			Description: db.Description.String,
			Stars:       int(db.Stars.Int32),
			Forks:       int(db.Forks.Int32),
			CreatedAt:   db.CreatedAt.Time.Format(time.RFC3339),
		}, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("database error: %w", err)
	}

	fetchRequest := kafka.FetchRequest{
		Owner: owner,
		Repo:  repo,
	}
	if err := u.kafkaProducer.SendFetchRequest(ctx, fetchRequest); err != nil {
		return nil, fmt.Errorf("failed to send fetch request to kafka: %w", err)
	}

	return nil, errors.New("try again later")
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

func (u *RepoInfo) GetSubscribedStats(ctx context.Context) ([]domain.Repository, error) {
	dbRepos, err := u.storage.ListAllStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all stats from storage: %w", err)
	}

	var result []domain.Repository
	for _, db := range dbRepos {
		result = append(result, domain.Repository{
			FullName:    db.FullName.String,
			Description: db.Description.String,
			Stars:       int(db.Stars.Int32),
			Forks:       int(db.Forks.Int32),
			CreatedAt:   db.CreatedAt.Time.Format(time.RFC3339),
		})
	}

	return result, nil
}

package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"repo-stat/collector/internal/domain"
	platformKafka "repo-stat/platform/kafka"
	"time"
)

type RepoInfo struct {
	ghClient  GitHubClient
	subClient SubscriptionLister
	log       *slog.Logger
}

func NewRepoInfo(ghClient GitHubClient, subClient SubscriptionLister, log *slog.Logger) *RepoInfo {
	return &RepoInfo{
		ghClient:  ghClient,
		subClient: subClient,
		log:       log,
	}
}

func (u *RepoInfo) GetRepoInfo(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	ghRepo, err := u.ghClient.GetRepoInfo(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("github api error: %w", err)
	}

	return &domain.Repository{
		FullName:    ghRepo.FullName,
		Description: ghRepo.Description,
		Stars:       ghRepo.Stars,
		Forks:       ghRepo.Forks,
		CreatedAt:   ghRepo.CreatedAt,
	}, nil
}

func (u *RepoInfo) GetSubscribedStats(ctx context.Context) ([]*domain.Repository, error) {
	subs, err := u.subClient.GetActiveSubscriptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active subscriptions: %w", err)
	}

	var repos []*domain.Repository

	for _, sub := range subs {
		repoData, err := u.ghClient.GetRepoInfo(ctx, sub.Owner, sub.Repo)
		if err != nil {
			u.log.Warn("failed to fetch repo info for subscription",
				"owner", sub.Owner, "repo", sub.Repo, "error", err)
			continue
		}
		repos = append(repos, repoData)
	}

	return repos, nil
}

func (u *RepoInfo) StartBackgroundUpdate(ctx context.Context, producer KafkaProducer) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			u.log.Info("starting background subscription update")
			subs, err := u.subClient.GetActiveSubscriptions(ctx)
			if err != nil {
				u.log.Error("failed to get subscriptions for background update", "error", err)
				continue
			}

			for _, sub := range subs {
				err := producer.SendFetchRequest(ctx, platformKafka.FetchRequest{
					Owner: sub.Owner,
					Repo:  sub.Repo,
				})
				if err != nil {
					u.log.Error("failed to send background update request", "repo", sub.Repo, "error", err)
				}
			}
		}
	}
}

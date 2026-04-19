package usecase

import (
	"context"
	"errors"
	"repo-stat/subscriber/internal/adapter/github"
	"repo-stat/subscriber/internal/adapter/storage"
	"repo-stat/subscriber/internal/domain"
)

type SubscriptionManager struct {
	repo     *storage.Repository
	ghClient *github.Client
}

func NewSubscriptionManager(repo *storage.Repository, ghClient *github.Client) *SubscriptionManager {
	return &SubscriptionManager{
		repo:     repo,
		ghClient: ghClient,
	}
}

func (sm *SubscriptionManager) AddSubscription(ctx context.Context, sub *domain.Subscription) error {
	if !sub.IsValid() {
		return errors.New("invalid subscription data")
	}

	err := sm.ghClient.CheckRepoExists(ctx, sub.Owner, sub.Repo)
	if err != nil {
		if err.Error() == "repository not found" {
			return errors.New("repository does not exist on GitHub")
		}
		return err
	}

	return sm.repo.AddSubscription(ctx, sub)
}

func (sm *SubscriptionManager) ListSubscriptions(ctx context.Context) ([]domain.Subscription, error) {
	return sm.repo.ListSubscriptions(ctx)
}

func (sm *SubscriptionManager) RemoveSubscription(ctx context.Context, owner, repo string) error {
	sub := &domain.Subscription{Owner: owner, Repo: repo}
	if !sub.IsValid() {
		return errors.New("invalid owner or repo")
	}
	return sm.repo.RemoveSubscription(ctx, owner, repo)
}

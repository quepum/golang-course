package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type SubscriptionManager struct {
	subClient SubscriberClient
}

func NewSubscriptionManager(subClient SubscriberClient) *SubscriptionManager {
	return &SubscriptionManager{subClient: subClient}
}

func (s *SubscriptionManager) AddSubscription(ctx context.Context, owner, repo string) error {
	return s.subClient.AddSubscription(ctx, owner, repo)
}

func (s *SubscriptionManager) ListSubscriptions(ctx context.Context) ([]domain.Subscription, error) {
	return s.subClient.ListSubscriptions(ctx)
}

func (s *SubscriptionManager) RemoveSubscription(ctx context.Context, owner, repo string) error {
	return s.subClient.RemoveSubscription(ctx, owner, repo)
}

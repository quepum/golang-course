package storage

import (
	"context"
	"repo-stat/subscriber/internal/domain"
)

type Repository struct {
	db *Queries
}

func NewRepository(db *Queries) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AddSubscription(ctx context.Context, sub *domain.Subscription) error {
	_, err := r.db.CreateSubscription(ctx, CreateSubscriptionParams{
		Owner: sub.Owner,
		Repo:  sub.Repo,
	})
	return err
}

func (r *Repository) ListSubscriptions(ctx context.Context) ([]domain.Subscription, error) {
	rows, err := r.db.ListSubscriptions(ctx)
	if err != nil {
		return nil, err
	}

	var subs []domain.Subscription
	for _, row := range rows {
		subs = append(subs, domain.Subscription{
			Owner: row.Owner,
			Repo:  row.Repo,
		})
	}
	return subs, nil
}

func (r *Repository) RemoveSubscription(ctx context.Context, owner, repo string) error {
	return r.db.DeleteSubscription(ctx, DeleteSubscriptionParams{
		Owner: owner,
		Repo:  repo,
	})
}

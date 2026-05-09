package subscriber

import (
	"context"
	"log/slog"
	"repo-stat/api/internal/domain"

	subscirberpb "repo-stat/proto/subscriber"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   subscirberpb.SubscriberClient
}

func (c *Client) GetName() string {
	return "subscriber"
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		log:  log,
		conn: conn,
		pb:   subscirberpb.NewSubscriberClient(conn),
	}, nil
}

func (c *Client) Ping(ctx context.Context) domain.PingStatus {
	_, err := c.pb.Ping(ctx, &subscirberpb.PingRequest{})
	if err != nil {
		c.log.Error("subscriber ping failed", "error", err)
		return domain.PingStatusDown
	}

	return domain.PingStatusUp
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) AddSubscription(ctx context.Context, owner, repo string) error {
	_, err := c.pb.AddSubscription(ctx, &subscirberpb.AddSubscriptionRequest{
		Owner: owner,
		Repo:  repo,
	})
	return err
}

func (c *Client) RemoveSubscription(ctx context.Context, owner, repo string) error {
	_, err := c.pb.RemoveSubscription(ctx, &subscirberpb.RemoveSubscriptionRequest{
		Owner: owner,
		Repo:  repo,
	})
	return err
}

func (c *Client) ListSubscriptions(ctx context.Context) ([]domain.Subscription, error) {
	resp, err := c.pb.ListSubscriptions(ctx, &subscirberpb.ListSubscriptionsRequest{})
	if err != nil {
		return nil, err
	}

	var subs []domain.Subscription
	for _, item := range resp.Subscriptions {
		subs = append(subs, domain.Subscription{
			Owner: item.Owner,
			Repo:  item.Repo,
		})
	}
	return subs, nil
}

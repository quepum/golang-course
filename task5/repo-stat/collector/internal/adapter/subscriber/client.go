package subscriber

import (
	"context"
	"log/slog"
	"repo-stat/collector/internal/domain"
	subscriberv1 "repo-stat/proto/subscriber"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   subscriberv1.SubscriberClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		log:  log,
		conn: conn,
		pb:   subscriberv1.NewSubscriberClient(conn),
	}, nil
}

func (c *Client) GetActiveSubscriptions(ctx context.Context) ([]*domain.Subscription, error) {
	resp, err := c.pb.ListSubscriptions(ctx, &subscriberv1.ListSubscriptionsRequest{})
	if err != nil {
		return nil, err
	}

	var subs []*domain.Subscription
	for _, s := range resp.Subscriptions {
		subs = append(subs, &domain.Subscription{
			Owner: s.Owner,
			Repo:  s.Repo,
		})
	}

	return subs, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

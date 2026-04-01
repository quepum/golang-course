package collector

import (
	"context"
	"log/slog"
	"repo-stat/api/internal/domain"
	collectorv1 "repo-stat/proto/collector"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   collectorv1.CollectorClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		log:  log,
		conn: conn,
		pb:   collectorv1.NewCollectorClient(conn),
	}, nil
}

func (c *Client) GetName() string {
	return "collector"
}

func (c *Client) Ping(ctx context.Context) domain.PingStatus {
	_, err := c.pb.Ping(ctx, &collectorv1.PingRequest{})
	if err != nil {
		c.log.Error("Collector ping failed", "error", err)
		return domain.PingStatusDown
	}
	return domain.PingStatusUp
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

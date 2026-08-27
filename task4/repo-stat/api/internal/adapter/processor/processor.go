package processor

import (
	"context"
	"log/slog"
	"repo-stat/api/internal/domain"
	processorv1 "repo-stat/proto/processor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   processorv1.ProcessorClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		log:  log,
		conn: conn,
		pb:   processorv1.NewProcessorClient(conn),
	}, nil
}

func (c *Client) GetName() string {
	return "processor"
}

func (c *Client) Ping(ctx context.Context) domain.PingStatus {
	_, err := c.pb.Ping(ctx, &processorv1.PingRequest{})
	if err != nil {
		c.log.Error("Processor ping failed", "error", err)
		return domain.PingStatusDown
	}
	return domain.PingStatusUp
}

func (c *Client) GetRepoInfo(ctx context.Context, url string) (*domain.Repository, error) {
	resp, err := c.pb.GetRepoInfo(ctx, &processorv1.RepoInfoRequest{
		Url: url,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Repository{
		FullName:    resp.FullName,
		Description: resp.Description,
		Stars:       resp.Stars,
		Forks:       resp.Forks,
		CreatedAt:   resp.CreatedAt,
	}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) GetSubscribedStats(ctx context.Context) ([]domain.Repository, error) {
	resp, err := c.pb.GetSubscribedStats(ctx, &processorv1.SubscribedStatsRequest{})
	if err != nil {
		return nil, err
	}

	var repos []domain.Repository
	for _, stat := range resp.Stats {
		repos = append(repos, domain.Repository{
			FullName:    stat.FullName,
			Description: stat.Description,
			Stars:       stat.Stars,
			Forks:       stat.Forks,
			CreatedAt:   stat.CreatedAt,
		})
	}
	return repos, nil
}

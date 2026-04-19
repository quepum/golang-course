package collector

import (
	"context"
	"log/slog"
	"repo-stat/processor/internal/domain"
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

func (c *Client) GetRepoInfo(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	resp, err := c.pb.GetRepoInfo(ctx, &collectorv1.RepoInfoRequest{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Repository{
		FullName:    resp.FullName,
		Description: resp.Description,
		Stars:       int(resp.Stars),
		Forks:       int(resp.Forks),
		CreatedAt:   resp.CreatedAt,
	}, nil
}

func (c *Client) GetSubscribedStats(ctx context.Context) ([]*domain.Repository, error) {
	resp, err := c.pb.GetSubscribedStats(ctx, &collectorv1.SubscribedStatsRequest{})
	if err != nil {
		return nil, err
	}

	var repos []*domain.Repository
	for _, stat := range resp.Stats {
		repos = append(repos, &domain.Repository{
			FullName:    stat.FullName,
			Description: stat.Description,
			Stars:       int(stat.Stars),
			Forks:       int(stat.Forks),
			CreatedAt:   stat.CreatedAt,
		})
	}
	return repos, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

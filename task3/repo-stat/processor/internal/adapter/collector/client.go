package collector

import (
	"context"
	"log/slog"
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

func (c *Client) GetRepoInfo(ctx context.Context, owner, repo string) (*collectorv1.RepoInfoResponse, error) {
	return c.pb.GetRepoInfo(ctx, &collectorv1.RepoInfoRequest{
		Owner: owner,
		Repo:  repo,
	})
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

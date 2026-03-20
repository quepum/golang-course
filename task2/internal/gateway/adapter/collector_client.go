package adapter

import (
	"context"
	"task2/domain"
	pb "task2/pkg/proto/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type collectorClient struct {
	client pb.RepoServiceClient
}

func NewCollectorClient(conn *grpc.ClientConn) domain.Repository {
	return &collectorClient{client: pb.NewRepoServiceClient(conn)}
}

func (c *collectorClient) GetRepoInfo(ctx context.Context, owner, repo string) (*domain.RepoInfo, error) {
	request := &pb.RepoRequest{Owner: owner, Repo: repo}
	response, err := c.client.GetRepoInfo(ctx, request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case 3:
				return nil, domain.ErrInvalidInput
			case 5:
				return nil, domain.ErrRepoNotFound
			default:
				return nil, err
			}
		}

		return nil, err
	}

	info := &domain.RepoInfo{
		Name:        response.Name,
		Description: response.Description,
		Stars:       int(response.Stars),
		Forks:       int(response.Forks),
		CreatedAt:   response.CreatedAt,
	}

	return info, nil
}

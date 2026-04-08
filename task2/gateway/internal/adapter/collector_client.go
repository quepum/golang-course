package adapter

import (
	"context"
	"task2/gateway/domain"
	pb "task2/pkg/proto/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
			case codes.InvalidArgument:
				return nil, domain.ErrInvalidInput
			case codes.NotFound:
				return nil, domain.ErrRepoNotFound
			default:
				return nil, err
			}
		}

		return nil, err
	}

	info := &domain.RepoInfo{
		Name:        response.GetName(),
		Description: response.GetDescription(),
		Stars:       int(response.GetStars()),
		Forks:       int(response.GetForks()),
		CreatedAt:   response.GetCreatedAt(),
	}

	return info, nil
}

package handler

import (
	"context"
	"errors"
	"task2/collector/domain"
	pb "task2/pkg/proto/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	pb.UnimplementedRepoServiceServer
	uc domain.UseCase
}

func NewGrpcServer(uc domain.UseCase) *GrpcServer {
	return &GrpcServer{uc: uc}
}

func (s *GrpcServer) GetRepoInfo(ctx context.Context, req *pb.RepoRequest) (*pb.RepoResponse, error) {
	info, err := s.uc.GetRepoInfo(ctx, req.GetOwner(), req.GetRepo())
	if err != nil {
		if errors.Is(err, domain.ErrRepoNotFound) {
			return nil, status.Error(codes.NotFound, "repository not found")
		}

		if errors.Is(err, domain.ErrInvalidInput) {
			return nil, status.Error(codes.InvalidArgument, "invalid input parameters")
		}

		return nil, status.Error(codes.Internal, "failed to fetch repository information")
	}

	response := &pb.RepoResponse{
		Name:        info.Name,
		Description: info.Description,
		Stars:       int32(info.Stars),
		Forks:       int32(info.Forks),
		CreatedAt:   info.CreatedAt,
	}

	return response, nil
}

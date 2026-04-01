package usecase

import (
	"context"
	"fmt"
	"repo-stat/collector/internal/adapter/github"
	collectorv1 "repo-stat/proto/collector"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RepoInfo struct {
	ghClient *github.Client
}

func NewRepoInfo(ghClient *github.Client) *RepoInfo {
	return &RepoInfo{ghClient: ghClient}
}

func (u *RepoInfo) Execute(ctx context.Context, owner, repo string) (*collectorv1.RepoInfoResponse, error) {
	ghRepo, err := u.ghClient.GetRepoInfo(ctx, owner, repo)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "repository not found" {
			return nil, status.Error(codes.NotFound, errMsg)
		}
		if errMsg == "rate limit exceeded" {
			return nil, status.Error(codes.ResourceExhausted, errMsg)
		}
		return nil, status.Error(codes.Internal, fmt.Errorf("github api error: %w", err).Error())
	}

	return &collectorv1.RepoInfoResponse{
		FullName:    ghRepo.FullName,
		Description: ghRepo.Description,
		Stars:       int32(ghRepo.Stars),
		Forks:       int32(ghRepo.Forks),
		CreatedAt:   ghRepo.CreatedAt,
	}, nil
}

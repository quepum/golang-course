package usecase

import (
	"context"
	"fmt"
	"repo-stat/collector/internal/adapter/github"
	"repo-stat/collector/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RepoInfo struct {
	ghClient *github.Client
}

func NewRepoInfo(ghClient *github.Client) *RepoInfo {
	return &RepoInfo{ghClient: ghClient}
}

func (u *RepoInfo) Execute(ctx context.Context, owner, repo string) (*domain.Repository, error) {
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

	return &domain.Repository{
		FullName:    ghRepo.FullName,
		Description: ghRepo.Description,
		Stars:       ghRepo.Stars,
		Forks:       ghRepo.Forks,
		CreatedAt:   ghRepo.CreatedAt,
	}, nil
}

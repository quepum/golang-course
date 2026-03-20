package usecase

import (
	"context"
	"task2/domain"
)

type repoUseCase struct {
	repo domain.Repository
}

func NewRepoUseCase(repo domain.Repository) domain.UseCase {
	return &repoUseCase{repo: repo}
}

func (r *repoUseCase) GetRepoInfo(ctx context.Context, owner, repo string) (*domain.RepoInfo, error) {
	if owner == "" || repo == "" {
		return nil, domain.ErrInvalidInput
	}

	return r.repo.GetRepoInfo(ctx, owner, repo)
}

package domain

import (
	"context"
	"errors"
)

type RepoInfo struct {
	Name        string
	Description string
	Stars       int
	Forks       int
	CreatedAt   string
}

type Repository interface {
	GetRepoInfo(ctx context.Context, owner, repo string) (*RepoInfo, error)
}

type UseCase interface {
	GetRepoInfo(ctx context.Context, owner, repo string) (*RepoInfo, error)
}

var (
	ErrRepoNotFound = errors.New("repo not found")
	ErrInvalidInput = errors.New("invalid input")
)

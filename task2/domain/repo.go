package domain

import (
	"context"
	"errors"
)

type RepoInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int    `json:"stars"`
	Forks       int    `json:"forks"`
	CreatedAt   string `json:"created_at"`
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

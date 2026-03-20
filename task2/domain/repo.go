package domain

import (
	"context"
	"errors"
)

type RepoInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	CreatedAt   string `json:"created_at"`
}

type Repository interface {
	GetRepoInfo(ctx context.Context, owner, repo string) (*RepoInfo, error)
}

type UseCase interface {
	GetRepoInfo(ctx context.Context, owner, repo string) (*RepoInfo, error)
}

var ErrRepoNotFound = errors.New("repo not found")

var ErrInvalidInput = errors.New("invalid input")

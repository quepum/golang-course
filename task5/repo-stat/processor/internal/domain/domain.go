package domain

import "errors"

type Repository struct {
	FullName    string
	Description string
	Stars       int
	Forks       int
	CreatedAt   string
}

var ErrDataNotReady = errors.New("repository data is not ready yet")

package storage

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	*Queries
}

func NewRepository(db *pgxpool.Pool) *Repo {
	return &Repo{
		Queries: New(db),
	}
}

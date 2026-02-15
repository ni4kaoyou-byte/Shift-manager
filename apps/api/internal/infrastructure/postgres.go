package infrastructure

import "github.com/jackc/pgx/v5/pgxpool"

type Postgres struct {
	Pool *pgxpool.Pool
}

func NewPostgres(pool *pgxpool.Pool) *Postgres {
	return &Postgres{Pool: pool}
}

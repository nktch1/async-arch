package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func initStorage(ctx context.Context) (*sqlx.DB, error) {
	pgDSN := os.Getenv("PG_DSN")
	if pgDSN == "" {
		pgDSN = "postgres://postgres:password@localhost:5433/postgres?sslmode=disable"
	}

	connection, err := sqlx.ConnectContext(ctx, "postgres", pgDSN)
	if err != nil {
		return nil, fmt.Errorf("create sql connect: %w", err)
	}

	return connection, nil
}

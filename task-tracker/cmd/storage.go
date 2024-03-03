package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx"
)

func initStorage(ctx context.Context) (*pgx.Conn, error) {
	pgDSN := os.Getenv("PG_DSN")
	if pgDSN == "" {
		pgDSN = "postgres://postgres:password@localhost:5433/postgres"
	}

	config, err := pgx.ParseConnectionString(pgDSN)
	if err != nil {
		return nil, fmt.Errorf("parse conn string: %w", err)
	}

	connection, err := pgx.Connect(config)
	if err != nil {
		return nil, fmt.Errorf("pg connect: %w", err)
	}

	if err = applyMigrations(ctx, connection); err != nil {
		return nil, fmt.Errorf("apply migrations: %w", err)
	}

	return connection, nil
}

func applyMigrations(ctx context.Context, connection *pgx.Conn) error {
	var migrations []string

	clientsCreateTable := `
	CREATE TABLE IF NOT EXISTS clients (
	  id     	 TEXT  NOT NULL,
	  secret 	 TEXT  NOT NULL,
	  domain 	 TEXT  NOT NULL,
	  CONSTRAINT clients_pkey PRIMARY KEY (id)
	)`

	tasksCreateTable := `
	CREATE TABLE IF NOT EXISTS tasks (
	  id 			bigserial PRIMARY KEY,
	  description   TEXT 	  NOT NULL,
	  is_open		bool 	  NOT NULL,
	  public_id		UUID 	  NOT NULL,
	  popug_id 		TEXT 	  REFERENCES clients (id)
	)`

	migrations = append(migrations, clientsCreateTable)
	migrations = append(migrations, tasksCreateTable)

	for _, migration := range migrations {
		_, err := connection.Exec(migration)
		if err != nil {
			return fmt.Errorf("apply migration: %w", err)
		}
	}

	return nil
}

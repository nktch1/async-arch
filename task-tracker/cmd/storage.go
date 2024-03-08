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

	if err = applyMigrations(ctx, connection); err != nil {
		return nil, fmt.Errorf("apply migrations: %w", err)
	}

	return connection, nil
}

func applyMigrations(ctx context.Context, connection *sqlx.DB) error {
	var migrations []string

	clientsCreateTable := `
	CREATE TABLE IF NOT EXISTS clients (
	  id     	 TEXT  NOT NULL,
	  secret 	 TEXT  NOT NULL,
	  domain 	 TEXT  NOT NULL,
	  CONSTRAINT clients_pkey PRIMARY KEY (id)
	)`

	//testUser := `INSERT INTO clients (id, secret, domain) VALUES ('1', 'secret', 'domain')`

	tasksCreateTable := `
	CREATE TABLE IF NOT EXISTS tasks (
	  id 			bigserial PRIMARY KEY,
	  description   TEXT 	  NOT NULL,
	  jira_id       TEXT 	  NOT NULL,
	  is_open		bool 	  NOT NULL,
	  public_id		UUID 	  NOT NULL,
	  popug_id 		TEXT 	  REFERENCES clients (id)
	)`

	migrations = append(migrations, clientsCreateTable)
	//migrations = append(migrations, testUser)
	migrations = append(migrations, tasksCreateTable)

	for _, migration := range migrations {
		_, err := connection.ExecContext(ctx, migration)
		if err != nil {
			return fmt.Errorf("apply migration: %w", err)
		}
	}

	return nil
}

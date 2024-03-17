package tasksrepo

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	taskmodel "github.com/nikitych1/awesome-task-exchange-system/accounting/internal/entity/task"
)

type DB struct {
	connection *sqlx.DB
}

func New(connection *sqlx.DB) DB {
	return DB{connection: connection}
}

func (d DB) AddTask(ctx context.Context, task taskmodel.Task) error {
	fmt.Println(task.PublicID)
	insertQuery := `
INSERT INTO tasks 
    (public_id, account_public_id, description, status, cost, reward) 
VALUES 
    ($1, $2, $3, $4, $5, $6)
ON CONFLICT (public_id) DO UPDATE
SET
    account_public_id = excluded.account_public_id, 
    status = excluded.status,
    updated_at = NOW()
`
	_, err := d.connection.ExecContext(
		ctx,
		insertQuery,
		task.PublicID,
		task.AccountPublicID,
		task.Description,
		task.Status,
		task.Cost,
		task.Reward,
	)
	if err != nil {
		return fmt.Errorf("add task repo: %w", err)
	}

	return nil
}

func (d DB) ShuffleTasks(ctx context.Context, task taskmodel.Task) error {
	updateQuery := `
UPDATE tasks
SET account_public_id = $1
WHERE public_id = $2`

	_, err := d.connection.ExecContext(
		ctx,
		updateQuery,
		task.AccountPublicID,
		task.PublicID,
	)
	if err != nil {
		return fmt.Errorf("shuffle task repo: %w", err)
	}

	return nil
}

func (d DB) CloseTask(ctx context.Context, task taskmodel.Task) error {
	updateQuery := `
UPDATE tasks
SET status = $1
WHERE public_id = $2`

	_, err := d.connection.ExecContext(
		ctx,
		updateQuery,
		task.Status,
		task.PublicID,
	)
	if err != nil {
		return fmt.Errorf("close task repo: %w", err)
	}

	return nil
}

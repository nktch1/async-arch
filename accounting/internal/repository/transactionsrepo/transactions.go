package transactionsrepo

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	taskmodel "github.com/nikitych1/awesome-task-exchange-system/accounting/internal/entity/task"
	"github.com/nikitych1/awesome-task-exchange-system/accounting/internal/entity/transaction"
)

type DB struct {
	connection *sqlx.DB
}

func New(connection *sqlx.DB) DB {
	return DB{connection: connection}
}

func (d DB) ListTransactions(ctx context.Context) ([]transaction.Transaction, error) {
	return nil, nil
}

func (d DB) ListTransactionsByAccount(ctx context.Context, accountPublicID uuid.UUID) ([]transaction.Transaction, error) {
	return nil, nil
}

func (d DB) ChargeMoney(ctx context.Context, task taskmodel.Task) error {
	insertQuery := `
INSERT INTO transactions
	(task_public_id, account_public_id, amount)
VALUES
    ($1, $2, $3)
`
	_, err := d.connection.ExecContext(
		ctx,
		insertQuery,
		task.PublicID,
		task.AccountPublicID,
		-task.Cost,
	)
	if err != nil {
		return fmt.Errorf("shuffle task repo: %w", err)
	}

	return nil
}

func (d DB) PayMoney(ctx context.Context, task taskmodel.Task) error {
	insertQuery := `
INSERT INTO transactions
	(task_public_id, account_public_id, amount)
SELECT public_id AS task_public_id, account_public_id, reward AS amount 
FROM tasks
WHERE public_id = $1
`
	_, err := d.connection.ExecContext(
		ctx,
		insertQuery,
		task.PublicID,
	)
	if err != nil {
		return fmt.Errorf("shuffle task repo: %w", err)
	}

	return nil
}

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
	selectQuery := `
	SELECT id, account_public_id, task_public_id, credit, debit, created_at, updated_at
	FROM transactions
`
	rows, err := d.connection.QueryxContext(ctx, selectQuery)
	if err != nil {
		return nil, fmt.Errorf("select transactions by account id repo: %w", err)
	}

	var transactions []transaction.Transaction

	for rows.Next() {
		var transactionToScan transactionEntity

		if err = rows.Scan(
			&transactionToScan.ID,
			&transactionToScan.AccountPublicID,
			&transactionToScan.TaskPublicID,
			&transactionToScan.Credit,
			&transactionToScan.Debit,
			&transactionToScan.CreatedAt,
			&transactionToScan.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}

		transaction, err := transactionEntityToTransaction(transactionToScan)
		if err != nil {
			return nil, fmt.Errorf("convert transaction: %w", err)
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (d DB) ListTransactionsByAccount(ctx context.Context, accountPublicID uuid.UUID) ([]transaction.Transaction, error) {
	selectQuery := `
	SELECT id, account_public_id, task_public_id, credit, debit, created_at, updated_at
	FROM transactions
	WHERE account_public_id = $1
`
	rows, err := d.connection.QueryxContext(ctx, selectQuery, accountPublicID)
	if err != nil {
		return nil, fmt.Errorf("select transactions by account id repo: %w", err)
	}

	var transactions []transaction.Transaction

	for rows.Next() {
		var transactionToScan transactionEntity

		if err = rows.Scan(
			&transactionToScan.ID,
			&transactionToScan.AccountPublicID,
			&transactionToScan.TaskPublicID,
			&transactionToScan.Credit,
			&transactionToScan.Debit,
			&transactionToScan.CreatedAt,
			&transactionToScan.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}

		transaction, err := transactionEntityToTransaction(transactionToScan)
		if err != nil {
			return nil, fmt.Errorf("convert transaction: %w", err)
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (d DB) ChargeMoney(ctx context.Context, task taskmodel.Task) error {
	insertQuery := `
INSERT INTO transactions
	(task_public_id, account_public_id, debit)
VALUES
    ($1, $2, $3)
`
	_, err := d.connection.ExecContext(
		ctx,
		insertQuery,
		task.PublicID,
		task.AccountPublicID,
		task.Cost,
	)
	if err != nil {
		return fmt.Errorf("shuffle task repo: %w", err)
	}

	return nil
}

func (d DB) PayMoney(ctx context.Context, task taskmodel.Task) error {
	insertQuery := `
INSERT INTO transactions
	(task_public_id, account_public_id, credit)
SELECT public_id AS task_public_id, account_public_id, reward AS credit 
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

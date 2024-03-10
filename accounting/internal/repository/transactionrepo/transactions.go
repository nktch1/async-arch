package transactionrepo

import (
	"context"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

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

package accounting

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/nikitych1/awesome-task-exchange-system/accounting/internal/entity/transaction"
)

type transactionRepository interface {
	ListTransactions(context.Context) ([]transaction.Transaction, error)
	ListTransactionsByAccount(context.Context, uuid.UUID) ([]transaction.Transaction, error)
}

type Accounting struct {
	transactionRepository
}

func New(transactionRepo transactionRepository) Accounting {
	return Accounting{transactionRepository: transactionRepo}
}

func (a Accounting) ListTransactions(context.Context) ([]transaction.Transaction, error) {
	return nil, nil
}

func (a Accounting) ListTransactionsByAccount(ctx context.Context, accountPublicID uuid.UUID) ([]transaction.Transaction, error) {
	return nil, nil
}

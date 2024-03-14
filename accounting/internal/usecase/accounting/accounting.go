package accounting

import (
	"context"
	"fmt"

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

type ListTransactionsResponse struct {
	Income float64 `json:"income"`

	Transactions []transaction.Transaction `json:"transactions"`
}

func (a Accounting) ListTransactions(ctx context.Context) (ListTransactionsResponse, error) {
	var response ListTransactionsResponse

	transactions, err := a.transactionRepository.ListTransactions(ctx)
	if err != nil {
		return ListTransactionsResponse{}, fmt.Errorf("read transactions: %w", err)
	}

	response.Transactions = transactions

	for _, transaction := range transactions {
		response.Income += transaction.Credit
		response.Income -= transaction.Debit
	}

	return response, nil
}

type ListTransactionsByAccountResponse struct {
	Income float64 `json:"income"`

	Transactions []transaction.Transaction `json:"transactions"`
}

func (a Accounting) ListTransactionsByAccount(ctx context.Context, accountPublicID uuid.UUID) (ListTransactionsByAccountResponse, error) {
	var response ListTransactionsByAccountResponse

	transactions, err := a.transactionRepository.ListTransactionsByAccount(ctx, accountPublicID)
	if err != nil {
		return ListTransactionsByAccountResponse{}, fmt.Errorf("read transactions: %w", err)
	}

	response.Transactions = transactions

	for _, transaction := range transactions {
		response.Income += transaction.Credit
		response.Income -= transaction.Debit
	}

	return response, nil
}

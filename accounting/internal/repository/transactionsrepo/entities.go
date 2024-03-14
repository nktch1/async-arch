package transactionsrepo

import (
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/nikitych1/awesome-task-exchange-system/accounting/internal/entity/transaction"
)

type transactionEntity struct {
	ID int `db:"id"`

	PublicID        string `db:"public_id"`
	AccountPublicID string `db:"account_public_id"`
	TaskPublicID    string `db:"task_public_id"`

	Credit float64 `db:"credit"`
	Debit  float64 `db:"debit"`

	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func transactionEntityToTransaction(t transactionEntity) (transaction.Transaction, error) {
	accountPublicIDUUID, err := uuid.FromString(t.AccountPublicID)
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("parse account public id uuid: %w", err)
	}

	taskPublicID, err := uuid.FromString(t.TaskPublicID)
	if err != nil {
		return transaction.Transaction{}, fmt.Errorf("parse task public id uuid: %w", err)
	}

	return transaction.Transaction{
		AccountPublicID: accountPublicIDUUID,
		TaskPublicID:    taskPublicID,

		Credit: t.Credit,
		Debit:  t.Debit,
	}, nil
}

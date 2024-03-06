package accountsrepo

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/entity/account"
)

type DB struct {
	connection *sqlx.DB
}

func New(connection *sqlx.DB) DB {
	return DB{connection: connection}
}

func (d DB) CreateAccount(ctx context.Context) error {
	return nil
}

func (d DB) ListAssigneeAccounts(ctx context.Context) ([]account.Account, error) {
	selectQuery := `
SELECT (id, public_id, name, email, role, created_at, updated_at)
FROM accounts
WHERE role in ('employee', 'manager')`

	rows, err := d.connection.QueryContext(ctx, selectQuery)
	if err != nil {
		return nil, fmt.Errorf("list accounts repo: %w", err)
	}

	accounts := make([]account.Account, 0)

	for rows.Next() {
		scannedAccount := accountEntity{}

		err = rows.Scan(
			&scannedAccount.ID,
			&scannedAccount.PublicID,
			&scannedAccount.Name, &scannedAccount.Email,
			&scannedAccount.Role, &scannedAccount.CreatedAt,
			&scannedAccount.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan account: %w", err)
		}

		task, err := accountEntityToAccount(scannedAccount)
		if err != nil {
			return nil, fmt.Errorf("convert task: %w", err)
		}

		accounts = append(accounts, task)
	}

	return accounts, nil
}

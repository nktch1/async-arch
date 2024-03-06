package accountsrepo

import (
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/entity/account"
)

type accountEntity struct {
	ID int `db:"id"`

	PublicID string `db:"public_id"`

	Name  string `db:"name"`
	Email string `db:"email"`

	Role account.Role `db:"role"`

	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func accountEntityToAccount(a accountEntity) (account.Account, error) {
	publicIDUUID, err := uuid.FromString(a.PublicID)
	if err != nil {
		return account.Account{}, fmt.Errorf("parse public id uuid: %w", err)
	}

	return account.Account{
		Name:  a.Name,
		Email: a.Email,
		Role:  a.Role,

		PublicID: publicIDUUID,
	}, nil
}

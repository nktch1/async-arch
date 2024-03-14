package transaction

import uuid "github.com/satori/go.uuid"

type Transaction struct {
	AccountPublicID uuid.UUID `json:"account_public_id"`
	TaskPublicID    uuid.UUID `json:"task_public_id"`

	Amount float64 `json:"amount"`
}

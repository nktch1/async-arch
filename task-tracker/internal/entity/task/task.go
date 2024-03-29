package task

import uuid "github.com/satori/go.uuid"

type Status string

const (
	New  Status = "new"
	Done Status = "done"
)

type Task struct {
	PublicID        uuid.UUID `json:"public_id"`
	AccountPublicID uuid.UUID `json:"account_public_id"`
	Description     string    `json:"description"`
	Status          Status    `json:"status"`
}

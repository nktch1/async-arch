package tasksrepo

import (
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/entity/task"
)

type taskEntity struct {
	ID int `db:"id"`

	PublicID        string `db:"public_id"`
	AccountPublicID string `db:"account_public_id"`
	Description     string `db:"description"`

	Status task.Status `db:"status"`

	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func taskEntityToTask(t taskEntity) (task.Task, error) {
	publicIDUUID, err := uuid.FromString(t.PublicID)
	if err != nil {
		return task.Task{}, fmt.Errorf("parse public id uuid: %w", err)
	}

	accountPublicIDUUID, err := uuid.FromString(t.AccountPublicID)
	if err != nil {
		return task.Task{}, fmt.Errorf("parse account public id uuid: %w", err)
	}

	return task.Task{
		PublicID:        publicIDUUID,
		AccountPublicID: accountPublicIDUUID,
		Description:     t.Description,
		Status:          t.Status,
	}, nil
}

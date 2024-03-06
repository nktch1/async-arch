package tasksrepo

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/entity/account"
	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/entity/task"
)

type DB struct {
	connection *sqlx.DB
}

func New(connection *sqlx.DB) DB {
	return DB{connection: connection}
}

func (d DB) ListTasks(ctx context.Context, accountPublicID uuid.UUID) ([]task.Task, error) {
	selectQuery := `
SELECT id, public_id, account_public_id, description, status, created_at, updated_at 
FROM tasks 
WHERE account_public_id = $1
`
	rows, err := d.connection.QueryContext(ctx, selectQuery, accountPublicID)
	if err != nil {
		return nil, fmt.Errorf("list tasks repo: %w", err)
	}

	tasks := make([]task.Task, 0)

	for rows.Next() {
		scannedTask := taskEntity{}

		err = rows.Scan(
			&scannedTask.ID, &scannedTask.PublicID,
			&scannedTask.AccountPublicID, &scannedTask.Description,
			&scannedTask.Status, &scannedTask.CreatedAt,
			&scannedTask.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}

		task, err := taskEntityToTask(scannedTask)
		if err != nil {
			return nil, fmt.Errorf("convert task: %w", err)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (d DB) AddTask(ctx context.Context, task task.Task) error {
	insertQuery := `
INSERT INTO tasks 
    (account_public_id, description, status) 
VALUES 
    ($1, $2, $3)`

	_, err := d.connection.ExecContext(
		ctx,
		insertQuery,
		task.AccountPublicID,
		task.Description,
		task.Status,
	)
	if err != nil {
		return fmt.Errorf("add task repo: %w", err)
	}

	return nil
}

func (d DB) ShuffleTasks(ctx context.Context, accounts []account.Account) error {
	selectQuery := `
SELECT id, public_id, account_public_id, description, status, created_at, updated_at 
FROM tasks 
WHERE status = 'new'`

	rows, err := d.connection.QueryContext(ctx, selectQuery)
	if err != nil {
		return fmt.Errorf("list tasks repo: %w", err)
	}

	openedTasks := make([]task.Task, 0)

	for rows.Next() {
		scannedTask := taskEntity{}

		err = rows.Scan(
			&scannedTask.ID, &scannedTask.PublicID,
			&scannedTask.AccountPublicID, &scannedTask.Description,
			&scannedTask.Status, &scannedTask.CreatedAt,
			&scannedTask.UpdatedAt,
		)

		if err != nil {
			return fmt.Errorf("scan task: %w", err)
		}

		task, err := taskEntityToTask(scannedTask)
		if err != nil {
			return fmt.Errorf("convert task: %w", err)
		}

		openedTasks = append(openedTasks, task)
	}

	getRandomAccount := func(accounts []account.Account) account.Account {
		minIdx, maxIdx := 0, len(accounts)-1
		return accounts[rand.Intn(maxIdx-minIdx)+minIdx]
	}

	updateQuery := `
UPDATE tasks
SET account_public_id = $1
WHERE public_id = $2`

	for _, openedTask := range openedTasks {
		randomAssigneeID := getRandomAccount(accounts).PublicID

		_, err = d.connection.ExecContext(
			ctx,
			updateQuery,
			randomAssigneeID,
			openedTask.PublicID,
		)
		if err != nil {
			return fmt.Errorf("close task repo: %w", err)
		}
	}

	return nil
}

func (d DB) CloseTask(ctx context.Context, taskPublicID uuid.UUID) error {
	updateQuery := `
UPDATE tasks 
SET status = 'done'
WHERE public_id = $1`

	_, err := d.connection.ExecContext(
		ctx,
		updateQuery,
		taskPublicID,
	)
	if err != nil {
		return fmt.Errorf("close task repo: %w", err)
	}

	return nil
}
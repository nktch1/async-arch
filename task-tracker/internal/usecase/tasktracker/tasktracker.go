package tasktracker

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/entity/account"
	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/entity/task"
)

type taskRepository interface {
	AddTask(context.Context, task.Task) (task.Task, error)
	ShuffleTasks(context.Context, []account.Account) ([]task.Task, error)
	CloseTask(context.Context, uuid.UUID) (task.Task, error)
	ListTasks(context.Context, uuid.UUID) ([]task.Task, error)
}

type taskProducer interface {
	ProduceAddedTaskEvent(context.Context, task.Task) error
	ProduceShuffledTasksEvents(context.Context, []task.Task) error
	ProduceClosedTaskEvent(context.Context, task.Task) error
}

type accountsRepository interface {
	ListAssigneeAccounts(context.Context) ([]account.Account, error)
}

type TaskTracker struct {
	accountsRepository
	taskRepository
	taskProducer
}

func New(taskRepo taskRepository, accRepo accountsRepository, taskProducer taskProducer) TaskTracker {
	return TaskTracker{taskRepository: taskRepo, accountsRepository: accRepo, taskProducer: taskProducer}
}

func (d TaskTracker) ListTasks(ctx context.Context, accountPublicID uuid.UUID) ([]task.Task, error) {
	return d.taskRepository.ListTasks(ctx, accountPublicID)
}

func (d TaskTracker) AddTask(ctx context.Context, task task.Task) error {
	affectedTask, err := d.taskRepository.AddTask(ctx, task)
	if err != nil {
		return fmt.Errorf("add task: %w", err)
	}

	if err = d.taskProducer.ProduceAddedTaskEvent(ctx, affectedTask); err != nil {
		return fmt.Errorf("send shuffled tasks event: %w", err)
	}

	return nil
}

func (d TaskTracker) ShuffleTasks(ctx context.Context) error {
	accounts, err := d.accountsRepository.ListAssigneeAccounts(ctx)
	if err != nil {
		return fmt.Errorf("list accounts: %w", err)
	}

	affectedTasks, err := d.taskRepository.ShuffleTasks(ctx, accounts)
	if err != nil {
		return fmt.Errorf("shuffle tasks: %w", err)
	}

	if err = d.taskProducer.ProduceShuffledTasksEvents(ctx, affectedTasks); err != nil {
		return fmt.Errorf("send shuffled tasks event: %w", err)
	}

	return nil
}

func (d TaskTracker) CloseTask(ctx context.Context, taskPublicID uuid.UUID) error {
	affectedTask, err := d.taskRepository.CloseTask(ctx, taskPublicID)
	if err != nil {
		return fmt.Errorf("close task: %w", err)
	}

	if err = d.taskProducer.ProduceClosedTaskEvent(ctx, affectedTask); err != nil {
		return fmt.Errorf("send closed task event: %w", err)
	}

	return nil
}

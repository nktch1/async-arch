package tasktracker

import (
	"context"
	"encoding/json"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"

	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/entity/account"
	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/entity/task"
)

type taskRepository interface {
	AddTask(context.Context, task.Task) error
	ShuffleTasks(context.Context, []account.Account) error
	CloseTask(context.Context, uuid.UUID) error
	ListTasks(context.Context, uuid.UUID) ([]task.Task, error)
}

type accountsRepository interface {
	ListAssigneeAccounts(context.Context) ([]account.Account, error)
}

type kafkaWriter interface {
	WriteMessages(...kafka.Message) (int, error)
}

type TaskTracker struct {
	accountsRepository
	taskRepository
	kafkaWriter
}

func New(taskRepo taskRepository, accRepo accountsRepository, writer kafkaWriter) TaskTracker {
	return TaskTracker{taskRepository: taskRepo, accountsRepository: accRepo, kafkaWriter: writer}
}

func (d TaskTracker) ListTasks(ctx context.Context, accountPublicID uuid.UUID) ([]task.Task, error) {
	return d.taskRepository.ListTasks(ctx, accountPublicID)
}

func (d TaskTracker) AddTask(ctx context.Context, task task.Task) error {
	if err := d.taskRepository.AddTask(ctx, task); err != nil {
		return fmt.Errorf("add task: %w", err)
	}

	if err := d.sendKafkaMessage(ctx, task); err != nil {
		return fmt.Errorf("send kafka event: %w", err)
	}

	return nil
}

func (d TaskTracker) ShuffleTasks(ctx context.Context) error {
	accounts, err := d.accountsRepository.ListAssigneeAccounts(ctx)
	if err != nil {
		return fmt.Errorf("list accounts: %w", err)
	}

	if err = d.taskRepository.ShuffleTasks(ctx, accounts); err != nil {
		return fmt.Errorf("shuffle tasks: %w", err)
	}

	return nil
}

func (d TaskTracker) CloseTask(ctx context.Context, taskPublicID uuid.UUID) error {
	if err := d.taskRepository.CloseTask(ctx, taskPublicID); err != nil {
		return fmt.Errorf("close task: %w", err)
	}

	return nil
}

func (d TaskTracker) sendKafkaMessage(ctx context.Context, object interface{}) error {
	// TODO separate Task and TaskEvent

	taskContent, err := json.Marshal(object)
	if err != nil {
		return fmt.Errorf("marshal task: %w", err)
	}

	messages := []kafka.Message{{Value: taskContent}}

	if _, err = d.kafkaWriter.WriteMessages(messages...); err != nil {
		return fmt.Errorf("write task event: %w", err)
	}

	return nil
}

package tasktracker

import (
	"context"

	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/entity/task"
)

type taskRepository interface {
	AddTask(context.Context, task.Task) error
	ShuffleTasks(context.Context) error
	CloseTask(context.Context) error
	ListTasks(context.Context) ([]task.Task, error)
}

type TaskTracker struct {
	repository taskRepository
}

func New(repo taskRepository) TaskTracker {
	return TaskTracker{repository: repo}
}

func (d TaskTracker) ListTasks(ctx context.Context) ([]task.Task, error) {
	return d.repository.ListTasks(ctx)
}

func (d TaskTracker) AddTask(ctx context.Context, t task.Task) error {
	return d.repository.AddTask(ctx, t)
}

func (d TaskTracker) ShuffleTasks(ctx context.Context) error {
	return nil
}

func (d TaskTracker) CloseTask(ctx context.Context) error {
	return nil
}

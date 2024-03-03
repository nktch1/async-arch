package tasktracker

import "github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/entity/task"

type taskRepository interface {
	AddTask(task.Task) error
	ShuffleTasks() error
	CloseTask() error
	ListTasks() ([]task.Task, error)
}

type TaskTracker struct {
	repository taskRepository
}

func New(repo taskRepository) TaskTracker {
	return TaskTracker{repository: repo}
}

func (d TaskTracker) ListTasks() ([]task.Task, error) {

	return d.repository.ListTasks()
}

func (d TaskTracker) AddTask(t task.Task) error {
	return d.repository.AddTask(t)
}

func (d TaskTracker) ShuffleTasks() error {
	return nil
}

func (d TaskTracker) CloseTask() error {
	return nil
}

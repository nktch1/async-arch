package taskworkfloweventproducer

import (
	"context"
	"fmt"

	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/pkg/kafka"

	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/entity/task"
)

const tasksWorkflowBusinessEventsTopic = "tasks-workflow"

type TaskWorkflowEventProducer struct {
	kafka.SRProducer
}

func New(srProducer kafka.SRProducer) TaskWorkflowEventProducer {
	return TaskWorkflowEventProducer{SRProducer: srProducer}
}

func (p TaskWorkflowEventProducer) ProduceAddedTaskEvent(ctx context.Context, task task.Task) error {
	_, err := p.SRProducer.ProduceMessage(taskToTaskAddedEvent(task), tasksWorkflowBusinessEventsTopic)
	return err
}

func (p TaskWorkflowEventProducer) ProduceShuffledTasksEvents(ctx context.Context, tasks []task.Task) error {
	for _, task := range tasks {
		_, err := p.SRProducer.ProduceMessage(taskToTaskShuffledEvent(task), tasksWorkflowBusinessEventsTopic)
		if err != nil {
			return fmt.Errorf("produce message: %w", err)
		}
	}

	return nil
}

func (p TaskWorkflowEventProducer) ProduceClosedTaskEvent(ctx context.Context, task task.Task) error {
	_, err := p.SRProducer.ProduceMessage(taskToTaskClosedEvent(task), tasksWorkflowBusinessEventsTopic)
	return err
}

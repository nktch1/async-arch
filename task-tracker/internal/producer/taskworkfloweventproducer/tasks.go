package taskworkfloweventproducer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"

	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/entity/task"
)

type kafkaWriter interface {
	WriteMessages(...kafka.Message) (int, error)
}

type TaskWorkflowEventProducer struct {
	kafkaWriter
}

func New(writer kafkaWriter) TaskWorkflowEventProducer {
	return TaskWorkflowEventProducer{kafkaWriter: writer}
}

func (p TaskWorkflowEventProducer) ProduceAddedTaskEvent(ctx context.Context, task task.Task) error {
	return p.sendKafkaMessage(ctx, taskToTaskAddedEvent(task))
}

func (p TaskWorkflowEventProducer) ProduceShuffledTasksEvents(ctx context.Context, tasks []task.Task) error {
	var shuffledEvents []interface{}

	for _, task := range tasks {
		shuffledEvents = append(shuffledEvents, taskToTaskShuffledEvent(task))
	}

	return p.sendKafkaMessage(ctx, shuffledEvents...)
}

func (p TaskWorkflowEventProducer) ProduceClosedTaskEvent(ctx context.Context, task task.Task) error {
	return p.sendKafkaMessage(ctx, taskToTaskClosedEvent(task))
}

func (p TaskWorkflowEventProducer) sendKafkaMessage(ctx context.Context, objects ...interface{}) error {
	if len(objects) == 0 {
		return nil
	}

	for _, object := range objects {
		taskContent, err := json.Marshal(object)
		if err != nil {
			return fmt.Errorf("marshal task: %w", err)
		}

		messages := []kafka.Message{{Value: taskContent}}

		if _, err = p.kafkaWriter.WriteMessages(messages...); err != nil {
			return fmt.Errorf("write task event: %w", err)
		}
	}

	return nil
}

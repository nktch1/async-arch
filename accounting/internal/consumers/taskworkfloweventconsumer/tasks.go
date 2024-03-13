package taskworkfloweventconsumer

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"google.golang.org/protobuf/proto"

	taskmodel "github.com/nikitych1/awesome-task-exchange-system/accounting/internal/entity/task"
	prototask "github.com/nikitych1/awesome-task-exchange-system/accounting/pkg/events/proto/task"
)

type taskRepository interface {
	AddTask(context.Context, taskmodel.Task) error
	ShuffleTasks(context.Context, taskmodel.Task) error
	CloseTask(context.Context, taskmodel.Task) error
}

type TaskWorkflowEventConsumer struct {
	taskRepository
}

func New(taskRepo taskRepository) TaskWorkflowEventConsumer {
	return TaskWorkflowEventConsumer{taskRepository: taskRepo}
}

func (t TaskWorkflowEventConsumer) Consume(ctx context.Context, event proto.Message) error {
	castedEvent, casted := event.(*prototask.TaskWorkflowEvent)
	if !casted {
		log.Printf("cast proto message to event type is failed: %+v", event)
		return nil
	}

	switch castedEvent.EventName {
	case prototask.TaskWorkflowEventType_TASK_ADDED_EVENT:
		return t.consumeAddedTaskEvent(ctx, castedEvent)
	case prototask.TaskWorkflowEventType_TASK_SHUFFLED_EVENT:
		return t.consumeShuffledTasksEvents(ctx, castedEvent)
	case prototask.TaskWorkflowEventType_TASK_CLOSED_EVENT:
		return t.consumeClosedTaskEvent(ctx, castedEvent)
	default:
		log.Printf("skip unknown event type: %s", castedEvent.String())
	}

	return nil
}

func (t TaskWorkflowEventConsumer) consumeAddedTaskEvent(ctx context.Context, event *prototask.TaskWorkflowEvent) error {
	getRandomCost := func() int {
		minIdx, maxIdx := 1, 10
		return rand.Intn(maxIdx-minIdx) + minIdx
	}

	cost := getRandomCost()

	model, err := taskEventToModelTask(event)
	if err != nil {
		return fmt.Errorf("convert event to model: %w", err)
	}

	model.Cost = cost

	return t.AddTask(ctx, model)
}

func (t TaskWorkflowEventConsumer) consumeShuffledTasksEvents(ctx context.Context, event *prototask.TaskWorkflowEvent) error {
	model, err := taskEventToModelTask(event)
	if err != nil {
		return fmt.Errorf("convert event to model: %w", err)
	}

	return t.taskRepository.ShuffleTasks(ctx, model)
}

func (t TaskWorkflowEventConsumer) consumeClosedTaskEvent(ctx context.Context, event *prototask.TaskWorkflowEvent) error {
	model, err := taskEventToModelTask(event)
	if err != nil {
		return fmt.Errorf("convert event to model: %w", err)
	}

	return t.taskRepository.CloseTask(ctx, model)
}

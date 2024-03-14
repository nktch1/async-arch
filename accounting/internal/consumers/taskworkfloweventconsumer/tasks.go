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

type transactionsRepository interface {
	ChargeMoney(context.Context, taskmodel.Task) error
	PayMoney(context.Context, taskmodel.Task) error
}

type tasksRepository interface {
	AddTask(context.Context, taskmodel.Task) error
	ShuffleTasks(context.Context, taskmodel.Task) error
	CloseTask(context.Context, taskmodel.Task) error
}

type TaskWorkflowEventConsumer struct {
	tasksRepository
	transactionsRepository
}

func New(taskRepo tasksRepository, transactionsRepo transactionsRepository) TaskWorkflowEventConsumer {
	return TaskWorkflowEventConsumer{tasksRepository: taskRepo, transactionsRepository: transactionsRepo}
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
		minIdx, maxIdx := 10, 20
		return rand.Intn(maxIdx-minIdx) + minIdx
	}

	task, err := taskEventToModelTask(event)
	if err != nil {
		return fmt.Errorf("convert event to model: %w", err)
	}

	task.Cost = getRandomCost()

	if err = t.tasksRepository.AddTask(ctx, task); err != nil {
		return fmt.Errorf("replicate add task: %w", err)
	}

	if err = t.transactionsRepository.ChargeMoney(ctx, task); err != nil {
		return fmt.Errorf("charge money: %w", err)
	}

	return nil
}

func (t TaskWorkflowEventConsumer) consumeShuffledTasksEvents(ctx context.Context, event *prototask.TaskWorkflowEvent) error {
	model, err := taskEventToModelTask(event)
	if err != nil {
		return fmt.Errorf("convert event to model: %w", err)
	}

	return t.tasksRepository.ShuffleTasks(ctx, model)
}

func (t TaskWorkflowEventConsumer) consumeClosedTaskEvent(ctx context.Context, event *prototask.TaskWorkflowEvent) error {
	task, err := taskEventToModelTask(event)
	if err != nil {
		return fmt.Errorf("convert event to model: %w", err)
	}

	if err = t.tasksRepository.CloseTask(ctx, task); err != nil {
		return fmt.Errorf("replicate close task: %w", err)
	}

	if err = t.transactionsRepository.PayMoney(ctx, task); err != nil {
		return fmt.Errorf("pay money: %w", err)
	}

	return nil
}

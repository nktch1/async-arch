package taskworkfloweventconsumer

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/protobuf/proto"

	prototask "github.com/nikitych1/awesome-task-exchange-system/accounting/pkg/events/proto/task"
)

type TaskWorkflowEventConsumer struct {
}

func New() TaskWorkflowEventConsumer {
	return TaskWorkflowEventConsumer{}
}

func (t TaskWorkflowEventConsumer) Consume(ctx context.Context, event proto.Message) error {
	castedEvent := event.(*prototask.TaskWorkflowEvent)

	switch castedEvent.EventName {
	case prototask.TaskWorkflowEventType_TASK_ADDED_EVENT:
		return t.consumeAddedTaskEvent(ctx)
	case prototask.TaskWorkflowEventType_TASK_SHUFFLED_EVENT:
		return t.consumeShuffledTasksEvents(ctx)
	case prototask.TaskWorkflowEventType_TASK_CLOSED_EVENT:
		return t.consumeClosedTaskEvent(ctx)
	default:
		log.Printf("skip unknown event type: %s", castedEvent.String())
	}

	return nil
}

func (t TaskWorkflowEventConsumer) consumeAddedTaskEvent(ctx context.Context) error {
	fmt.Println("added")
	return nil
}

func (t TaskWorkflowEventConsumer) consumeShuffledTasksEvents(ctx context.Context) error {
	fmt.Println("shuffled")
	return nil
}

func (t TaskWorkflowEventConsumer) consumeClosedTaskEvent(ctx context.Context) error {
	fmt.Println("closed")
	return nil
}

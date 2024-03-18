package taskworkfloweventproducer

import (
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/entity/task"
	prototask "github.com/nikitych1/awesome-task-exchange-system/task-tracker/pkg/events/proto/task"
)

func createBaseEvent(eventName prototask.TaskWorkflowEventType) *prototask.TaskWorkflowEvent {
	const (
		eventVersion = "1"
		producerName = "task-tracker"
	)

	return &prototask.TaskWorkflowEvent{
		EventName: eventName,
		Producer:  producerName,

		EventVersion: eventVersion,

		EventID:   uuid.NewV1().String(),
		EventTime: time.Now().String(),
	}
}

func taskToProtoTask(t task.Task) *prototask.Task {
	var status prototask.TaskStatus

	switch t.Status {
	case task.NewTaskStatus:
		status = prototask.TaskStatus_NEW_TASK_STATUS
	case task.DoneTaskStatus:
		status = prototask.TaskStatus_DONE_TASK_STATUS
	}

	return &prototask.Task{
		PublicID:        t.PublicID.String(),
		AccountPublicID: t.AccountPublicID.String(),
		Description:     t.Description,
		Status:          status,
	}
}

func taskToTaskAddedEvent(task task.Task) *prototask.TaskWorkflowEvent {
	event := createBaseEvent(prototask.TaskWorkflowEventType_TASK_ADDED_EVENT)
	event.TaskPayload = taskToProtoTask(task)
	return event
}

func taskToTaskShuffledEvent(task task.Task) *prototask.TaskWorkflowEvent {
	event := createBaseEvent(prototask.TaskWorkflowEventType_TASK_SHUFFLED_EVENT)
	event.TaskPayload = taskToProtoTask(task)
	return event
}

func taskToTaskClosedEvent(task task.Task) *prototask.TaskWorkflowEvent {
	event := createBaseEvent(prototask.TaskWorkflowEventType_TASK_CLOSED_EVENT)
	event.TaskPayload = taskToProtoTask(task)
	return event
}

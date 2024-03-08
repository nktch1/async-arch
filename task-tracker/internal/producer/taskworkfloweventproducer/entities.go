package taskworkfloweventproducer

import (
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/entity/task"
)

type EventType string

const (
	TaskAddedEventType    = "task_added_event"
	TaskClosedEventType   = "task_closed_event"
	TaskShuffledEventType = "task_shuffled_event"
)

type baseEvent struct {
	EventID      string `json:"event_id"`
	EventVersion string `json:"event_version"`
	EventTime    string `json:"event_time"`
	Producer     string `json:"event_producer_name"`

	EventName EventType `json:"event_name"`
}

type taskAddedEvent struct {
	baseEvent
	Data interface{}
}

type taskClosedEvent struct {
	baseEvent
	Data interface{}
}

type tasksShuffledEvent struct {
	baseEvent
	Data interface{}
}

func createBaseEvent(eventName EventType) baseEvent {
	const (
		eventVersion = "1"
		producerName = "task-tracker"
	)

	return baseEvent{
		EventName: eventName,
		Producer:  producerName,

		EventVersion: eventVersion,

		EventID:   uuid.NewV1().String(),
		EventTime: time.Now().String(),
	}
}

func taskToTaskAddedEvent(task task.Task) taskAddedEvent {
	return taskAddedEvent{baseEvent: createBaseEvent(TaskAddedEventType), Data: task}
}

func taskToTaskShuffledEvent(task task.Task) taskClosedEvent {
	return taskClosedEvent{baseEvent: createBaseEvent(TaskShuffledEventType), Data: task}
}

func taskToTaskClosedEvent(task task.Task) tasksShuffledEvent {
	return tasksShuffledEvent{baseEvent: createBaseEvent(TaskClosedEventType), Data: task}
}

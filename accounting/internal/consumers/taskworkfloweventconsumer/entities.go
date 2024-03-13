package taskworkfloweventconsumer

import (
	"fmt"

	uuid "github.com/satori/go.uuid"

	taskmodel "github.com/nikitych1/awesome-task-exchange-system/accounting/internal/entity/task"
	prototask "github.com/nikitych1/awesome-task-exchange-system/accounting/pkg/events/proto/task"
)

func taskEventToModelTask(event *prototask.TaskWorkflowEvent) (taskmodel.Task, error) {
	publicIDUUID, err := uuid.FromString(event.TaskPayload.PublicID)
	if err != nil {
		return taskmodel.Task{}, fmt.Errorf("parse public id uuid: %w", err)
	}

	accountPublicIDUUID, err := uuid.FromString(event.TaskPayload.AccountPublicID)
	if err != nil {
		return taskmodel.Task{}, fmt.Errorf("parse account public id uuid: %w", err)
	}

	var status taskmodel.Status

	switch event.TaskPayload.Status {
	case prototask.TaskStatus_NEW_TASK_STATUS:
		status = taskmodel.NewTaskStatus
	case prototask.TaskStatus_DONE_TASK_STATUS:
		status = taskmodel.DoneTaskStatus
	}

	return taskmodel.Task{
		Status: status,

		PublicID: publicIDUUID,

		Description:     event.TaskPayload.Description,
		AccountPublicID: accountPublicIDUUID,
	}, nil
}

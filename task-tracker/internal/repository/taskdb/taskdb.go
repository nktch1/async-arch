package taskdb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/entity/task"
	"github.com/segmentio/kafka-go"
)

type DB struct {
	connection  *pgx.Conn
	kafkaWriter *kafka.Writer
}

func New(connection *pgx.Conn, writer *kafka.Writer) DB {
	return DB{connection: connection, kafkaWriter: writer}
}

func (d DB) ListTasks() ([]task.Task, error) {
	rows, err := d.connection.Query("SELECT id, jira_id, description, is_open, popug_id FROM tasks where popug_id = $1")
	if err != nil {
		return nil, fmt.Errorf("list tasks repo: %w", err)
	}

	tasks := make([]task.Task, 0)

	for rows.Next() {
		var taskToScan taskEntity

		err = rows.Scan(&taskToScan.ID, &taskToScan.JiraID, &taskToScan.Description, &taskToScan.IsOpen, &taskToScan.PopugID)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		tasks = append(tasks, task.Task{
			Description: taskToScan.Description,
			JiraID:      taskToScan.JiraID,
			IsOpen:      taskToScan.IsOpen,
			PopugID:     taskToScan.PopugID,
		})
	}

	return tasks, nil
}

func (d DB) AddTask(task task.Task) error {
	taskUUID := uuid.New()

	_, err := d.connection.Exec(
		"INSERT INTO tasks (description, jira_id, is_open, popug_id, public_id) VALUES ($1, $2, $3, $4, $5)",
		task.Description,
		task.JiraID,
		task.IsOpen,
		task.PopugID,
		taskUUID,
	)
	if err != nil {
		return fmt.Errorf("add task repo: %w", err)
	}

	// TODO separate Task and TaskEvent

	taskContent, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("marsshal task: %w", err)
	}

	message := kafka.Message{Topic: "tasks-stream", Value: taskContent}

	if err = d.kafkaWriter.WriteMessages(context.TODO(), message); err != nil {
		return fmt.Errorf("write task event: %w", err)
	}

	return nil
}

func (d DB) ShuffleTasks() error {
	return nil
}

func (d DB) CloseTask() error {
	return nil
}

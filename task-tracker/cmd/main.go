package main

import (
	"context"
	"fmt"
	"log"

	handler "github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/gateway/openapi/tasktracker"
	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/producer/taskworkfloweventproducer"
	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/repository/accountsrepo"
	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/repository/tasksrepo"
	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/usecase/tasktracker"
)

func main() {
	ctx := context.Background()

	pgConnection, err := initStorage(ctx)
	if err != nil {
		log.Fatalf("init storage: %s", err.Error())
	}

	kafkaProducer, err := initKafka(ctx)
	defer kafkaProducer.Close()

	if err != nil {
		log.Fatalf("init kafka producer: %s", err.Error())
	}

	tasksRepository := tasksrepo.New(pgConnection)
	accountsRepository := accountsrepo.New(pgConnection)
	taskWorkflowEventProducer := taskworkfloweventproducer.New(kafkaProducer)
	taskTrackerService := tasktracker.New(tasksRepository, accountsRepository, taskWorkflowEventProducer)
	taskTrackerHandler := handler.New(taskTrackerService)

	fmt.Println("task tracker service is listening on :8080...")

	listenAndServe(ctx, taskTrackerHandler)

	fmt.Println("task-tracker service finished...")
}

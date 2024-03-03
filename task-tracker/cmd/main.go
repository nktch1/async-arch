package main

import (
	"context"
	"fmt"
	"log"

	handler "github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/gateways/openapi/tasktracker"
	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/repository/taskdb"
	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/internal/usecase/tasktracker"
)

func main() {
	ctx := context.Background()

	pgConnection, err := initStorage(ctx)
	if err != nil {
		log.Fatalf("init storage: %s", err.Error())
	}

	//kafkaReader, err := initKafkaReader(ctx)
	//if err != nil {
	//	log.Fatalf("init kafka reader: %s", err.Error())
	//}

	kafkaConnection, err := initKafka(ctx)
	if err != nil {
		log.Fatalf("init kafka writer: %s", err.Error())
	}

	taskTrackerRepository := taskdb.New(pgConnection, kafkaConnection)
	taskTrackerService := tasktracker.New(taskTrackerRepository)
	taskTrackerHandler := handler.New(taskTrackerService)

	listenAndServe(ctx, taskTrackerHandler)

	fmt.Println("task-tracker finished...")
}

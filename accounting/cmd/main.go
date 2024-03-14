package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nikitych1/awesome-task-exchange-system/accounting/internal/consumer/taskworkfloweventconsumer"
	handler "github.com/nikitych1/awesome-task-exchange-system/accounting/internal/gateway/openapi/accounting"
	"github.com/nikitych1/awesome-task-exchange-system/accounting/internal/repository/tasksrepo"
	"github.com/nikitych1/awesome-task-exchange-system/accounting/internal/repository/transactionsrepo"
	"github.com/nikitych1/awesome-task-exchange-system/accounting/internal/usecase/accounting"
)

func main() {
	ctx := context.Background()

	pgConnection, err := initStorage(ctx)
	if err != nil {
		log.Fatalf("init storage: %s", err.Error())
	}

	kafkaProducer, kafkaConsumer, err := initKafka(ctx)
	defer kafkaProducer.Close()
	defer kafkaConsumer.Close()

	if err != nil {
		log.Fatalf("init kafka producer: %s", err.Error())
	}

	transactionsRepository := transactionsrepo.New(pgConnection)
	tasksRepository := tasksrepo.New(pgConnection)
	accountingService := accounting.New(transactionsRepository)
	accountingHandler := handler.New(accountingService)
	taskWorkflowEventConsumer := taskworkfloweventconsumer.New(tasksRepository, transactionsRepository)

	go consumeAndServe(ctx, kafkaConsumer, taskWorkflowEventConsumer)

	fmt.Println("accounting service is listening on :8081...")

	listenAndServe(ctx, accountingHandler)

	fmt.Println("accounting service finished...")
}

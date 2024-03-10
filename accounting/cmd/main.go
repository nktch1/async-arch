package main

import (
	"context"
	"fmt"
	"log"

	handler "github.com/nikitych1/awesome-task-exchange-system/accounting/internal/gateway/openapi/accounting"
	"github.com/nikitych1/awesome-task-exchange-system/accounting/internal/repository/transactionrepo"
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

	transactionRepository := transactionrepo.New(pgConnection)
	accountingService := accounting.New(transactionRepository)
	accountingHandler := handler.New(accountingService)

	fmt.Println("accounting service is listening on :8081...")

	listenAndServe(ctx, accountingHandler)

	fmt.Println("accounting service finished...")
}

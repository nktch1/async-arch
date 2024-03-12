package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"

	"github.com/nikitych1/awesome-task-exchange-system/accounting/internal/consumers/taskworkfloweventconsumer"
	prototask "github.com/nikitych1/awesome-task-exchange-system/accounting/pkg/events/proto/task"
	"github.com/nikitych1/awesome-task-exchange-system/accounting/pkg/kafka"
)

func listenAndServe(ctx context.Context, handler *mux.Router) {
	httpServer := http.Server{Addr: ":8081", Handler: handler}

	idleConnectionsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Fatalf("HTTP server shutdown error: %v", err)
		}

		close(idleConnectionsClosed)
	}()

	if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server error: %v", err)
	}

	<-idleConnectionsClosed
}

func consumeAndServe(ctx context.Context, baseConsumer kafka.SRConsumer, taskWorkflowConsumer taskworkfloweventconsumer.TaskWorkflowEventConsumer) {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := baseConsumer.Close(); err != nil {
			log.Fatalf("HTTP server shutdown error: %v", err)
		}
	}()

	const tasksWorkflowBusinessEventsTopic = "task-workflow"

	messageType := (&prototask.TaskWorkflowEvent{}).ProtoReflect().Type()

	if err := baseConsumer.Run(ctx, messageType, tasksWorkflowBusinessEventsTopic, taskWorkflowConsumer); err != nil {
		log.Fatalf("kafka consumer error: %v", err)
	}
}

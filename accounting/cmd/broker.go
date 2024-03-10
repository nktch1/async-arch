package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nikitych1/awesome-task-exchange-system/accounting/pkg/kafka"
)

func initKafka(ctx context.Context) (kafka.SRProducer, kafka.SRConsumer, error) {
	kafkaDSN := os.Getenv("KAFKA_DSN")
	if kafkaDSN == "" {
		kafkaDSN = "localhost:9092"
	}

	schemaRegistryURL := os.Getenv("SCHEMA_REGISTRY_URL")
	if kafkaDSN == "" {
		kafkaDSN = "localhost:8085"
	}

	producer, err := kafka.NewProducer(kafkaDSN, schemaRegistryURL)
	if err != nil {
		return nil, nil, fmt.Errorf("create kafka producer: %w", err)
	}

	consumer, err := kafka.NewConsumer(kafkaDSN, schemaRegistryURL)
	if err != nil {
		return nil, nil, fmt.Errorf("create kafka consumer: %w", err)
	}

	return producer, consumer, nil
}

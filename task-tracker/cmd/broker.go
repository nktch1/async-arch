package main

import (
	"context"
	"os"

	"github.com/segmentio/kafka-go"
)

func initKafkaReader(ctx context.Context) (*kafka.Reader, error) {
	kafkaDSN := os.Getenv("KAFKA_DSN")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaDSN},
		Topic:   "accounts-stream",
	})

	return reader, nil
}

func initKafkaWriter(ctx context.Context) (*kafka.Writer, error) {
	kafkaDSN := os.Getenv("KAFKA_DSN")

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaDSN},
		Topic:   "tasks-stream",
	})

	return writer, nil
}

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/segmentio/kafka-go"
)

const tasksWorkflowBusinessEventsTopic = "tasks-workflow"

func initKafka(ctx context.Context) (*kafka.Conn, error) {
	kafkaDSN := os.Getenv("KAFKA_DSN")
	if kafkaDSN == "" {
		kafkaDSN = "localhost:9092"
	}

	conn, err := kafka.DialLeader(ctx, "tcp", kafkaDSN, tasksWorkflowBusinessEventsTopic, 0)
	if err != nil {
		return nil, fmt.Errorf("connect to kafka: %w", err)
	}

	if err = prepareTopics(ctx, conn); err != nil {
		return nil, fmt.Errorf("prepare kafka topics: %w", err)
	}

	return conn, nil
}

func prepareTopics(ctx context.Context, controllerConn *kafka.Conn) error {
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             tasksWorkflowBusinessEventsTopic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	if err := controllerConn.CreateTopics(topicConfigs...); err != nil {
		return fmt.Errorf("create kafka controller connection: %w", err)
	}

	return nil
}

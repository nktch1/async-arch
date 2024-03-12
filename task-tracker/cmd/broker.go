package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nikitych1/awesome-task-exchange-system/task-tracker/pkg/kafka"
)

const tasksWorkflowBusinessEventsTopic = "tasks-workflow"

func initKafka(ctx context.Context) (kafka.SRProducer, error) {
	kafkaDSN := os.Getenv("KAFKA_DSN")
	if kafkaDSN == "" {
		kafkaDSN = "localhost"
	}

	schemaRegistryURL := os.Getenv("SCHEMA_REGISTRY_URL")
	if schemaRegistryURL == "" {
		schemaRegistryURL = "http://localhost:8085"
	}

	producer, err := kafka.NewProducer(kafkaDSN, schemaRegistryURL)
	if err != nil {
		return nil, fmt.Errorf("create kafka producer: %w", err)
	}

	//producer, err := kafka.NewProducer(kafkaURL, schemaRegistryURL)

	//conn, err := kafka.DialLeader(ctx, "tcp", kafkaDSN, tasksWorkflowBusinessEventsTopic, 0)
	//if err != nil {
	//	return nil, fmt.Errorf("connect to kafka: %w", err)
	//}

	//if err = prepareTopics(ctx, conn); err != nil {
	//	return nil, fmt.Errorf("prepare kafka topics: %w", err)
	//}

	return producer, nil
}

//
//func prepareTopics(ctx context.Context, controllerConn *kafka.Conn) error {
//	topicConfigs := []kafka.TopicConfig{
//		{
//			Topic:             tasksWorkflowBusinessEventsTopic,
//			NumPartitions:     1,
//			ReplicationFactor: 1,
//		},
//	}
//
//	if err := controllerConn.CreateTopics(topicConfigs...); err != nil {
//		return fmt.Errorf("create kafka controller connection: %w", err)
//	}
//
//	return nil
//}

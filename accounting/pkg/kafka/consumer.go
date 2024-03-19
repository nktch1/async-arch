package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/protobuf"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	consumerGroupID       = "test-consumer"
	defaultSessionTimeout = 6000
	noTimeout             = -1
)

// SRConsumer interface
type SRConsumer interface {
	Run(ctx context.Context, messageType protoreflect.MessageType, topic string, handler eventsConsumer) error
	Close() error
}

type eventsConsumer interface {
	Consume(ctx context.Context, msg proto.Message) error
}

type srConsumer struct {
	consumer     *kafka.Consumer
	deserializer *protobuf.Deserializer
}

// NewConsumer returns new consumer with schema registry
func NewConsumer(kafkaURL, srURL string) (SRConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  kafkaURL,
		"group.id":           consumerGroupID,
		"session.timeout.ms": defaultSessionTimeout,
		"enable.auto.commit": false,
	})
	if err != nil {
		return nil, err
	}

	sr, err := schemaregistry.NewClient(schemaregistry.NewConfig(srURL))
	if err != nil {
		return nil, err
	}

	d, err := protobuf.NewDeserializer(sr, serde.ValueSerde, protobuf.NewDeserializerConfig())
	if err != nil {
		return nil, err
	}
	return &srConsumer{
		consumer:     c,
		deserializer: d,
	}, nil
}

// RegisterMessage add simpleHandler and register schema in SR
func (c *srConsumer) RegisterMessage(messageType protoreflect.MessageType) error {
	return nil
}

// Run consumer
func (c *srConsumer) Run(ctx context.Context, messageType protoreflect.MessageType, topic string, handler eventsConsumer) error {
	if err := c.consumer.SubscribeTopics([]string{topic}, nil); err != nil {
		return fmt.Errorf("subscribe topic: %w", err)
	}

	if err := c.deserializer.ProtoRegistry.RegisterMessage(messageType); err != nil {
		return fmt.Errorf("register message: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := c.handleMessage(ctx, topic, handler); err != nil {
				log.Printf("handle message error: %s", err.Error())
				continue
			}
		}
	}
}

func (c *srConsumer) handleMessage(ctx context.Context, topic string, handler eventsConsumer) error {
	kafkaMsg, err := c.consumer.ReadMessage(noTimeout)
	if err != nil {
		return fmt.Errorf("read message: %w", err)
	}

	msg, err := c.deserializer.Deserialize(topic, kafkaMsg.Value)
	if err != nil {
		return fmt.Errorf("deserialize: %w", err)
	}

	if err = handler.Consume(ctx, msg.(proto.Message)); err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	if _, err = c.consumer.CommitMessage(kafkaMsg); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

// Close all connections
func (c *srConsumer) Close() error {
	if err := c.consumer.Close(); err != nil {
		return fmt.Errorf("close: %w", err)
	}

	c.deserializer.Close()

	return nil
}

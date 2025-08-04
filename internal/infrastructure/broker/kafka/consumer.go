package kafka

import (
	"context"
	"log"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type MessageHandler func(message *kafka.Message) error

type Consumer struct {
	consumer *kafka.Consumer
}

func NewConsumer(config *KafkaConfig) *Consumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": strings.Join(config.Brokers, ","),
		"group.id":          config.GroupID,
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	return &Consumer{consumer: consumer}
}

func (c *Consumer) Subscribe(topic string, handler MessageHandler) error {
	return c.consumer.Subscribe(topic, nil)
}

func (c *Consumer) ConsumeLoop(ctx context.Context, handler MessageHandler) error {
	for {
		select {
		case <-ctx.Done():
			c.consumer.Close()
			return ctx.Err()
		default:
			msg, err := c.consumer.ReadMessage(-1)
			if err == nil {
				if err := handler(msg); err != nil {
					return err
				}
			} else {
				log.Printf("Error reading message: %v", err)
			}
		}

	}
}

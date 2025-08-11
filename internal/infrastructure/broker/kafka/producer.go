package kafka

import (
	"context"
	"log"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(brokers []string) *Producer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": strings.Join(brokers, ","),
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	return &Producer{producer: producer}
}

func (p *Producer) SendMessage(ctx context.Context, topic string, message []byte) error {
	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: message,
	}, nil)
}

func (p *Producer) Close() {
	p.producer.Close()
}

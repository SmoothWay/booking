package kafka

import (
	"context"
	"log"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(brokers []string, clientID string) *Producer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": strings.Join(brokers, ","),
		"client.id":         clientID,
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	return &Producer{producer: producer}
}

func (p *Producer) SendMessage(ctx context.Context, topic string, message []byte) error {
	log.Printf("Sending message to topic: %s, message: %s", topic, string(message))

	deliveryChan := make(chan kafka.Event, 1)
	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: message,
	}, deliveryChan)

	if err != nil {
		log.Printf("Failed to produce message: %v", err)
		return err
	}

	// Wait for delivery report
	select {
	case e := <-deliveryChan:
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Printf("Failed to deliver message: %v", ev.TopicPartition.Error)
				return ev.TopicPartition.Error
			}
			log.Printf("Message delivered to topic %s [%d] at offset %v",
				*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
		default:
			log.Printf("Ignored event: %v", ev)
		}
	case <-ctx.Done():
		log.Printf("Context cancelled while waiting for delivery")
		return ctx.Err()
	}
	p.producer.Flush(1000)
	return nil
}

func (p *Producer) Close() {
	p.producer.Close()
}

package kafka

import (
	"context"
	"log"
	"strings"

	"github.com/SmoothWay/booking/internal/config"
	"github.com/SmoothWay/booking/internal/domain"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type MessageHandler func(message *kafka.Message) error

type Consumer struct {
	consumer *kafka.Consumer
}

func newConsumer(config *KafkaConfig) *Consumer {
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

func InitKafkaConsumer(cfg *config.Config, bookingUsecase domain.BookingUsecase, userUsecase domain.UserUsecase, unitUsecase domain.UnitUsecase) *Consumer {
	// Initialize Kafka configuration
	kafkaConfig := &KafkaConfig{
		Brokers: cfg.Kafka.Brokers,
		GroupID: cfg.Kafka.GroupID,
	}

	consumer := newConsumer(kafkaConfig)

	// Set up topic subscriptions with handlers
	topics := map[string]func(msg *kafka.Message) error{
		"booking.created":   BookingCreateHandler(bookingUsecase),
		"booking.updated":   BookingUpdateHandler(bookingUsecase),
		"booking.cancelled": BookingCancelHandler(bookingUsecase),
		"user.created":      UserCreateHandler(userUsecase),
		"user.updated":      UserUpdateHandler(userUsecase),
		"unit.created":      UnitCreateHandler(unitUsecase),
	}

	// Subscribe to all topics
	for topic, handler := range topics {
		if err := consumer.Subscribe(topic, handler); err != nil {
			log.Printf("Failed to subscribe to topic %s: %v", topic, err)
		}
	}

	return consumer
}

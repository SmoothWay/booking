package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/SmoothWay/booking/internal/domain/events"
	"github.com/SmoothWay/booking/internal/infrastructure/broker/kafka"
)

func main() {
	// Test Kafka producer
	brokers := []string{"localhost:9092"}
	groupID := "booking"

	producer := kafka.NewProducer(brokers, groupID)
	defer producer.Close()

	// Create a test user event
	userEvent := events.UserCreatedEvent{
		ID:          "test-user-123",
		Name:        "Test User",
		Email:       "test@test.com",
		PhoneNumber: "1234567890",
		Role:        "user",
		Password:    "password123",
		CreatedAt:   time.Now(),
	}

	// Marshal the event
	eventJSON, err := json.Marshal(userEvent)
	if err != nil {
		log.Fatalf("Failed to marshal event: %v", err)
	}

	// Send the message
	ctx := context.Background()
	log.Printf("Sending test message to topic 'user.created': %s", string(eventJSON))

	err = producer.SendMessage(ctx, "user.created", eventJSON)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	log.Println("Test message sent successfully!")

	// Wait a bit to ensure message is sent
	time.Sleep(2 * time.Second)
}

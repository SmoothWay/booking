package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/SmoothWay/booking/internal/domain/events"
	"github.com/SmoothWay/booking/internal/infrastructure/broker/kafka"
)

func main() {
	fmt.Println("=== Kafka Retry & DLT Test Script ===")

	// Test producer
	brokers := []string{"localhost:9092"}
	groupID := "booking"

	fmt.Printf("Creating producer with brokers: %v\n", brokers)
	producer := kafka.NewProducer(brokers, groupID)
	defer producer.Close()

	ctx := context.Background()

	// Test 1: Successful message (should be acknowledged immediately)
	fmt.Println("\n--- Test 1: Successful Message ---")
	successEvent := events.UserCreatedEvent{
		ID:          "success-user-" + fmt.Sprintf("%d", time.Now().Unix()),
		Name:        "Success User",
		Email:       "success@example.com",
		PhoneNumber: "1234567890",
		Role:        "user",
		Password:    "password123",
		CreatedAt:   time.Now(),
	}

	successJSON, _ := json.Marshal(successEvent)
	fmt.Printf("Sending successful message: %s\n", string(successJSON))
	err := producer.SendMessage(ctx, "user.created", successJSON)
	if err != nil {
		log.Printf("Failed to send successful message: %v", err)
	}

	// Test 2: Message that will fail and be retried (test@test.com)
	fmt.Println("\n--- Test 2: Message That Will Fail and Be Retried ---")
	failEvent := events.UserCreatedEvent{
		ID:          "fail-user-" + fmt.Sprintf("%d", time.Now().Unix()),
		Name:        "Fail User",
		Email:       "test@test.com", // This will cause failure
		PhoneNumber: "1234567890",
		Role:        "user",
		Password:    "password123",
		CreatedAt:   time.Now(),
	}

	failJSON, _ := json.Marshal(failEvent)
	fmt.Printf("Sending failing message: %s\n", string(failJSON))
	err = producer.SendMessage(ctx, "user.created", failJSON)
	if err != nil {
		log.Printf("Failed to send failing message: %v", err)
	}

	// Test 3: Invalid JSON message (will fail parsing)
	fmt.Println("\n--- Test 3: Invalid JSON Message ---")
	invalidJSON := []byte(`{"invalid": "json", "missing": "required fields"}`)
	fmt.Printf("Sending invalid JSON: %s\n", string(invalidJSON))
	err = producer.SendMessage(ctx, "user.created", invalidJSON)
	if err != nil {
		log.Printf("Failed to send invalid JSON: %v", err)
	}

	fmt.Println("\n=== Test Messages Sent ===")
	fmt.Println("1. Success message: Should be acknowledged immediately")
	fmt.Println("2. Failing message (test@test.com): Will be retried 3 times, then moved to DLT")
	fmt.Println("3. Invalid JSON: Will fail parsing and be retried")
	fmt.Println("\nCheck your application logs to see the retry and DLT behavior!")
	fmt.Println("\nExpected behavior:")
	fmt.Println("- Success message: 'Message acknowledged'")
	fmt.Println("- Failing message: 'Handler error' -> retry -> 'Max retries exceeded' -> 'Sending to DLT'")
	fmt.Println("- Invalid JSON: 'Failed to unmarshal event' -> retry -> eventually to DLT")
}

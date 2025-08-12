package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/SmoothWay/booking/internal/domain"
	"github.com/SmoothWay/booking/internal/domain/events"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// BookingCreateHandler handles booking creation events
func BookingCreateHandler(uc domain.BookingUsecase) func(msg *kafka.Message) error {
	return func(msg *kafka.Message) error {
		var event events.BookingCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal event: %w", err)
		}

		return uc.CreateBooking(context.Background(), event)
	}
}

// BookingUpdateHandler handles booking update events
func BookingUpdateHandler(uc domain.BookingUsecase) func(msg *kafka.Message) error {
	return func(msg *kafka.Message) error {
		var event events.BookingUpdatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal event: %w", err)
		}

		return uc.UpdateBooking(context.Background(), event)
	}
}

// BookingCancelHandler handles booking cancellation events
func BookingCancelHandler(uc domain.BookingUsecase) func(msg *kafka.Message) error {
	return func(msg *kafka.Message) error {
		var event events.BookingCancelledEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal event: %w", err)
		}

		return uc.CancelBooking(context.Background(), event)
	}
}

// UserCreateHandler handles user creation events
func UserCreateHandler(uc domain.UserUsecase) func(msg *kafka.Message) error {
	log.Println("UserCreateHandler handler setuped")
	return func(msg *kafka.Message) error {
		var event events.UserCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal event: %w", err)
		}

		log.Printf("UserCreateHandler: processing user creation for %s", event.Email)
		log.Printf("UserCreateHandler: successfully send mail to user for %s", event.Email)
		return nil
	}
}

// UserCreateHandlerWithAck handles user creation events with acknowledgment control
func UserCreateHandlerWithAck(uc domain.UserUsecase) func(msg *kafka.Message) (bool, error) {
	log.Println("UserCreateHandlerWithAck handler setuped")
	return func(msg *kafka.Message) (bool, error) {
		var event events.UserCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			return false, fmt.Errorf("failed to unmarshal event: %w", err)
		}
		if event.Email == "test@test.com" {
			return false, fmt.Errorf("test user")
		}
		log.Printf("UserCreateHandlerWithAck: processing user creation for %s", event.Email)

		log.Printf("UserCreateHandlerWithAck: successfully send mail to user for %s", event.Email)
		// Return true to acknowledge the message (processing successful)
		return true, nil
	}
}

// UserUpdateHandler handles user update events
func UserUpdateHandler(uc domain.UserUsecase) func(msg *kafka.Message) error {
	return func(msg *kafka.Message) error {
		var event events.UserUpdatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal event: %w", err)
		}

		return uc.UpdateUser(context.Background(), event)
	}
}

// UnitCreateHandler handles unit creation events
func UnitCreateHandler(uc domain.UnitUsecase) func(msg *kafka.Message) error {
	return func(msg *kafka.Message) error {
		var event events.UnitCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal event: %w", err)
		}

		return uc.CreateUnit(context.Background(), event)
	}
}

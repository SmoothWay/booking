package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SmoothWay/booking/internal/domain"
	"github.com/SmoothWay/booking/internal/domain/events"
	"github.com/confluentinc/confluent-kafka-go/kafka"
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
	return func(msg *kafka.Message) error {
		var event events.UserCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal event: %w", err)
		}

		return uc.CreateUser(context.Background(), event)
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

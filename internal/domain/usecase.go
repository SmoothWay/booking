package domain

import (
	"context"

	"github.com/SmoothWay/booking/internal/domain/events"
)

// BookingUsecase defines the business logic interface for booking operations
type BookingUsecase interface {
	GetBookings(ctx context.Context, page int, pageSize int) ([]Booking, error)
	GetBookingByID(ctx context.Context, id string) (*Booking, error)
	CreateBooking(ctx context.Context, event events.BookingCreatedEvent) error
	UpdateBooking(ctx context.Context, event events.BookingUpdatedEvent) error
	CancelBooking(ctx context.Context, event events.BookingCancelledEvent) error
	GetBookingsByUserID(ctx context.Context, userID string, limit, offset int) ([]Booking, error)
}

// UserUsecase defines the business logic interface for user operations
type UserUsecase interface {
	GetUsers(ctx context.Context, page int, pageSize int) ([]User, error)
	GetUserByID(ctx context.Context, id string) (User, error)
	CreateUser(ctx context.Context, event events.UserCreatedEvent) error
	UpdateUser(ctx context.Context, event events.UserUpdatedEvent) error
}

// UnitUsecase defines the business logic interface for unit operations
type UnitUsecase interface {
	GetUnits(ctx context.Context, page int, pageSize int) ([]Unit, error)
	GetUnitByID(ctx context.Context, id string) (Unit, error)
	CreateUnit(ctx context.Context, event events.UnitCreatedEvent) error
}

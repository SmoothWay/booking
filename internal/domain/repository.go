package domain

import (
	"context"
	"time"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id string) error
	GetUsers(ctx context.Context, page int, pageSize int) ([]*User, error)
}

type UnitRepository interface {
	GetUnitByID(ctx context.Context, id string) (*Unit, error)
	CreateUnit(ctx context.Context, unit *Unit) error
	UpdateUnit(ctx context.Context, unit *Unit) error
	DeleteUnit(ctx context.Context, id string) error
	GetUnits(ctx context.Context, page int, pageSize int) ([]*Unit, error)
}

type BookingRepository interface {
	GetBookingByID(ctx context.Context, id string) (*Booking, error)
	CreateBooking(ctx context.Context, booking *Booking) error
	UpdateBooking(ctx context.Context, booking *Booking) error
	DeleteBooking(ctx context.Context, id string) error
	GetBookingsByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*Booking, error)
	GetBookingsByUnitID(ctx context.Context, unitID string, limit, offset int) ([]*Booking, error)
	GetBookings(ctx context.Context, page int, pageSize int) ([]*Booking, error)
	GetBookingsByUserID(ctx context.Context, userID string, limit, offset int) ([]*Booking, error)
}

type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Increment(ctx context.Context, key string) (int64, error)
}

package events

import "time"

// BookingCreatedEvent represents a booking creation event
type BookingCreatedEvent struct {
	ID        string    `json:"id"`
	UnitID    string    `json:"unit_id"`
	UserID    string    `json:"user_id"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
}

// BookingUpdatedEvent represents a booking update event
type BookingUpdatedEvent struct {
	ID        string    `json:"id"`
	UnitID    string    `json:"unit_id"`
	UserID    string    `json:"user_id"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BookingCancelledEvent represents a booking cancellation event
type BookingCancelledEvent struct {
	ID          string    `json:"id"`
	UnitID      string    `json:"unit_id"`
	UserID      string    `json:"user_id"`
	CancelledAt time.Time `json:"cancelled_at"`
	Reason      string    `json:"reason,omitempty"`
}

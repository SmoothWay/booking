package events

import "time"

// UnitCreatedEvent represents a unit creation event
type UnitCreatedEvent struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

// UnitUpdatedEvent represents a unit update event
type UnitUpdatedEvent struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UnitDeletedEvent represents a unit deletion event
type UnitDeletedEvent struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}

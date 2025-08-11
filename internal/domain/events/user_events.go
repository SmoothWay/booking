package events

import "time"

// UserCreatedEvent represents a user creation event
type UserCreatedEvent struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserUpdatedEvent represents a user update event
type UserUpdatedEvent struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Role        string    `json:"role"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserDeletedEvent represents a user deletion event
type UserDeletedEvent struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}

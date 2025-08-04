package domain

import "time"

type Service struct {
	ID          string
	Name        string
	Price       float64
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

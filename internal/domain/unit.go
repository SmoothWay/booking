package domain

import "time"

type Unit struct {
	ID          string
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

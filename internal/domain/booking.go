package domain

import "time"

type Booking struct {
	ID        string
	UnitID    string
	UserID    string
	Status    string
	Type      string
	StartDate time.Time
	EndDate   time.Time
}

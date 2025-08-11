package handler

import "time"

type Booking struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	UnitID    string    `json:"unit_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type GetBookingsResponse struct {
	Pageable
	Bookings []Booking `json:"bookings"`
}

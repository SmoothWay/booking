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

type CreateBookingRequest struct {
	UserID    string    `json:"user_id" validate:"required"`
	UnitID    string    `json:"unit_id" validate:"required"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
}

type CreateBookingResponse struct {
	ID string `json:"id"`
}

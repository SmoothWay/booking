package usecase

import (
	"context"

	"github.com/SmoothWay/booking/internal/domain"
	"github.com/SmoothWay/booking/internal/domain/events"
	"github.com/SmoothWay/booking/internal/infrastructure/broker/kafka"
)

type BookingUsecase struct {
	kafkaProducer *kafka.Producer
	bookingRepo   domain.BookingRepository
}

func NewBookingUsecase(kafkaProducer *kafka.Producer, bookingRepo domain.BookingRepository) *BookingUsecase {
	return &BookingUsecase{
		kafkaProducer: kafkaProducer,
		bookingRepo:   bookingRepo,
	}
}

func (uc *BookingUsecase) GetBookings(ctx context.Context, page int, pageSize int) ([]domain.Booking, error) {
	bookings, err := uc.bookingRepo.GetBookings(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	var bookingsDomain []domain.Booking
	for _, booking := range bookings {
		bookingsDomain = append(bookingsDomain, *booking)
	}
	return bookingsDomain, nil
}

func (uc *BookingUsecase) GetBookingByID(ctx context.Context, id string) (*domain.Booking, error) {
	booking, err := uc.bookingRepo.GetBookingByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return booking, nil
}

// CreateBooking handles booking creation from Kafka events
func (uc *BookingUsecase) CreateBooking(ctx context.Context, event events.BookingCreatedEvent) error {
	booking := &domain.Booking{
		ID:        event.ID,
		UnitID:    event.UnitID,
		UserID:    event.UserID,
		Status:    event.Status,
		Type:      event.Type,
		StartDate: event.StartDate,
		EndDate:   event.EndDate,
	}

	return uc.bookingRepo.CreateBooking(ctx, booking)
}

// UpdateBooking handles booking updates from Kafka events
func (uc *BookingUsecase) UpdateBooking(ctx context.Context, event events.BookingUpdatedEvent) error {
	booking := &domain.Booking{
		ID:        event.ID,
		UnitID:    event.UnitID,
		UserID:    event.UserID,
		Status:    event.Status,
		Type:      event.Type,
		StartDate: event.StartDate,
		EndDate:   event.EndDate,
	}

	return uc.bookingRepo.UpdateBooking(ctx, booking)
}

// CancelBooking handles booking cancellation from Kafka events
func (uc *BookingUsecase) CancelBooking(ctx context.Context, event events.BookingCancelledEvent) error {
	// Get the existing booking
	booking, err := uc.bookingRepo.GetBookingByID(ctx, event.ID)
	if err != nil {
		return err
	}

	// Update status to cancelled
	booking.Status = "cancelled"

	return uc.bookingRepo.UpdateBooking(ctx, booking)
}

func (uc *BookingUsecase) GetBookingsByUserID(ctx context.Context, userID string, limit, offset int) ([]domain.Booking, error) {
	bookings, err := uc.bookingRepo.GetBookingsByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var bookingsDomain []domain.Booking
	for _, booking := range bookings {
		bookingsDomain = append(bookingsDomain, *booking)
	}
	return bookingsDomain, nil
}

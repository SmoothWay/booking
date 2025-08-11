package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/SmoothWay/booking/internal/domain"
)

type bookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *bookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) GetBookingByID(ctx context.Context, id string) (*domain.Booking, error) {
	query := `SELECT id, unit_id, user_id, status, type, start_date, end_date FROM bookings WHERE id = $1`
	var booking domain.Booking
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&booking.ID, &booking.UnitID, &booking.UserID, &booking.Status,
		&booking.Type, &booking.StartDate, &booking.EndDate,
	)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) CreateBooking(ctx context.Context, booking *domain.Booking) error {
	query := `INSERT INTO bookings (id, unit_id, user_id, status, type, start_date, end_date) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
		booking.ID, booking.UnitID, booking.UserID, booking.Status,
		booking.Type, booking.StartDate, booking.EndDate,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *bookingRepository) UpdateBooking(ctx context.Context, booking *domain.Booking) error {
	query := `UPDATE bookings SET unit_id = $1, user_id = $2, status = $3, type = $4, start_date = $5, end_date = $6 WHERE id = $7`
	_, err := r.db.ExecContext(ctx, query,
		booking.UnitID, booking.UserID, booking.Status, booking.Type,
		booking.StartDate, booking.EndDate, booking.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *bookingRepository) DeleteBooking(ctx context.Context, id string) error {
	query := `DELETE FROM bookings WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *bookingRepository) GetBookingsByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*domain.Booking, error) {
	query := `SELECT id, unit_id, user_id, status, type, start_date, end_date FROM bookings WHERE start_date >= $1 AND end_date <= $2 LIMIT $3 OFFSET $4`
	rows, err := r.db.QueryContext(ctx, query, startDate, endDate, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*domain.Booking
	for rows.Next() {
		var booking domain.Booking
		err := rows.Scan(&booking.ID, &booking.UnitID, &booking.UserID, &booking.Status, &booking.Type, &booking.StartDate, &booking.EndDate)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, &booking)
	}
	return bookings, nil
}

func (r *bookingRepository) GetBookings(ctx context.Context, page int, pageSize int) ([]*domain.Booking, error) {
	query := `SELECT id, unit_id, user_id, status, type, start_date, end_date FROM bookings LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*domain.Booking
	for rows.Next() {
		var booking domain.Booking
		err := rows.Scan(&booking.ID, &booking.UnitID, &booking.UserID, &booking.Status, &booking.Type, &booking.StartDate, &booking.EndDate)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, &booking)
	}
	return bookings, nil
}

func (r *bookingRepository) GetBookingsByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Booking, error) {
	query := `SELECT id, unit_id, user_id, status, type, start_date, end_date FROM bookings WHERE user_id = $1 LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*domain.Booking
	for rows.Next() {
		var booking domain.Booking
		err := rows.Scan(&booking.ID, &booking.UnitID, &booking.UserID, &booking.Status, &booking.Type, &booking.StartDate, &booking.EndDate)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, &booking)
	}
	return bookings, nil
}

func (r *bookingRepository) GetBookingsByUnitID(ctx context.Context, unitID string, limit, offset int) ([]*domain.Booking, error) {
	query := `SELECT id, unit_id, user_id, status, type, start_date, end_date FROM bookings WHERE unit_id = $1 LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, unitID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*domain.Booking
	for rows.Next() {
		var booking domain.Booking
		err := rows.Scan(&booking.ID, &booking.UnitID, &booking.UserID, &booking.Status, &booking.Type, &booking.StartDate, &booking.EndDate)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, &booking)
	}
	return bookings, nil
}

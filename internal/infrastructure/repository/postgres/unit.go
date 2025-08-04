package postgres

import (
	"context"
	"database/sql"

	"github.com/SmoothWay/booking/internal/domain"
)

type unitRepository struct {
	db    *sql.DB
	cache *domain.CacheRepository
}

func NewUnitRepository(db *sql.DB, cache *domain.CacheRepository) *unitRepository {
	return &unitRepository{db: db, cache: cache}
}

func (r *unitRepository) GetUnitByID(ctx context.Context, id string) (*domain.Unit, error) {
	query := `SELECT id, name, description, price, created_at, updated_at, deleted_at FROM units WHERE id = $1`
	var unit domain.Unit
	err := r.db.QueryRowContext(ctx, query, id).Scan(&unit.ID, &unit.Name, &unit.Description, &unit.Price, &unit.CreatedAt, &unit.UpdatedAt, &unit.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *unitRepository) CreateUnit(ctx context.Context, unit *domain.Unit) error {
	query := `INSERT INTO units (id, name, description, price, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query, unit.ID, unit.Name, unit.Description, unit.Price, unit.CreatedAt, unit.UpdatedAt, unit.DeletedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *unitRepository) UpdateUnit(ctx context.Context, unit *domain.Unit) error {
	query := `UPDATE units SET name = $1, description = $2, price = $3, updated_at = $4, deleted_at = $5 WHERE id = $6`
	_, err := r.db.ExecContext(ctx, query, unit.Name, unit.Description, unit.Price, unit.UpdatedAt, unit.DeletedAt, unit.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *unitRepository) DeleteUnit(ctx context.Context, id string) error {
	query := `UPDATE units SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

package postgres

import (
	"context"
	"database/sql"

	"github.com/SmoothWay/booking/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUsers(ctx context.Context, page int, pageSize int) ([]*domain.User, error) {
	query := `SELECT id, name, email, phone_number, password, role, created_at, updated_at, deleted_at FROM users LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.PhoneNumber, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, name, email, phone_number, password, role, created_at, updated_at, deleted_at FROM users WHERE id = $1`
	var user domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.PhoneNumber, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, name, email, phone_number, password, role, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.PhoneNumber, user.Password, user.Role, user.CreatedAt, user.UpdatedAt, user.DeletedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET name = $1, email = $2, phone_number = $3, password = $4, role = $5, updated_at = $6, deleted_at = $7 WHERE id = $8`
	_, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.PhoneNumber, user.Password, user.Role, user.UpdatedAt, user.DeletedAt, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

package postgres

import (
	"database/sql"
	"fmt"

	"github.com/SmoothWay/booking/internal/config"
	"github.com/golang-migrate/migrate/v4"
)

func NewPostgres(cfg config.DB) (*sql.DB, error) {
	dbUrl := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	if err := runMigrations(dbUrl, cfg.MigrationDir); err != nil {
		return nil, err
	}
	return db, nil
}

func runMigrations(dbUrl string, migrationDir string) error {
	m, err := migrate.New("file://"+migrationDir, dbUrl)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

package postgres

import (
	"database/sql"
	"fmt"

	"github.com/SmoothWay/booking/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func NewPostgres(cfg config.DB) (*sql.DB, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
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

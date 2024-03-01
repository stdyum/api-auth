package config

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type DatabaseConfig struct {
	URL string `env:"URL"`
}

func ConnectToDatabase(config DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.URL)
	if err != nil {
		return nil, err
	}

	if err = runMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	goose.SetTableName("_migrations")
	return goose.Up(db, "scheme/migrations")
}

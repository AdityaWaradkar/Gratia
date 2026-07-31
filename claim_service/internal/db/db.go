package db

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/adityawaradkar/gratia/claim_service/internal/config"
)

const (
	maxOpenConns    = 25
	maxIdleConns    = 25
	connMaxLifetime = 5 * time.Minute
	connMaxIdleTime = 2 * time.Minute
)

// New creates and configures a PostgreSQL connection pool
func New() (*sqlx.DB, error) {

	db, err := sqlx.Connect(
		"postgres",
		config.AppConfig.DatabaseURL,
	)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(
		maxOpenConns,
	)

	db.SetMaxIdleConns(
		maxIdleConns,
	)

	db.SetConnMaxLifetime(
		connMaxLifetime,
	)

	db.SetConnMaxIdleTime(
		connMaxIdleTime,
	)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
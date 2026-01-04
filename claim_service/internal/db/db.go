package db

import (
	"errors"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	maxOpenConns    = 25
	maxIdleConns    = 25
	connMaxLifetime = 5 * time.Minute
)

func New() (*sqlx.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, errors.New("DATABASE_URL not set")
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

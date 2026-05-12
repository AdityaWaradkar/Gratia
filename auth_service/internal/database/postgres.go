package database

import (
	"context"
	"log"

	"github.com/adityawaradkar/gratia/auth_service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresConnection(cfg *config.Config) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	err = pool.Ping(context.Background())

	if err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	log.Println("Connected to Neon PostgreSQL")

	return pool
}
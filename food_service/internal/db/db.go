package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect initializes PostgreSQL connection pool
func Connect(databaseURL string) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to create db pool: %v", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("database connected successfully")
	return pool
}

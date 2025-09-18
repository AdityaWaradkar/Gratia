package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/adityawaradkar/gratia/auth_service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DBPool *pgxpool.Pool

// Connect initializes the Postgres connection pool
func Connect(ctx context.Context) *pgxpool.Pool {
	if DBPool != nil {
		return DBPool
	}

	cfg, err := pgxpool.ParseConfig(config.GetDatabaseURL())
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v", err)
	}

	// Optional: customize pool settings
	cfg.MaxConns = 20
	cfg.MinConns = 2
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.HealthCheckPeriod = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	DBPool = pool
	fmt.Println("Connected to Postgres successfully")
	return DBPool
}

// Close closes the DB connection pool
func Close() {
	if DBPool != nil {
		DBPool.Close()
	}
}

// Optional: helper to run migrations
func RunMigration(ctx context.Context, sql string) error {
	if DBPool == nil {
		return fmt.Errorf("DB pool not initialized")
	}

	_, err := DBPool.Exec(ctx, sql)
	return err
}

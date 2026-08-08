package db

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
)

// Connect initializes and returns a highly concurrent PostgreSQL connection pool
func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
    config, err := pgxpool.ParseConfig(databaseURL)
    if err != nil {
        return nil, fmt.Errorf("failed to parse database configuration: %w", err)
    }

    pool, err := pgxpool.NewWithConfig(ctx, config)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    if err := pool.Ping(ctx); err != nil {
        return nil, fmt.Errorf("database ping failed: %w", err)
    }

    return pool, nil
}
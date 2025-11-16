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

func Connect(ctx context.Context) *pgxpool.Pool {
	if DBPool != nil {
		return DBPool
	}

	poolConfig, err := pgxpool.ParseConfig(config.GetDatabaseURL())
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v", err)
	}

	poolConfig.MaxConns = 20
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	DBPool = pool
	fmt.Println("Connected to Postgres successfully")
	return DBPool
}

func Close() {
	if DBPool != nil {
		DBPool.Close()
	}
}

func RunMigration(ctx context.Context, sql string) error {
	if DBPool == nil {
		return fmt.Errorf("DB pool not initialized")
	}
	_, err := DBPool.Exec(ctx, sql)
	return err
}

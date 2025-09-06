package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func ConnectDB() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set in environment")
	}

	// Connect to Neon Postgres
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Unable to create DB connection pool: ", err)
	}

	// Test connection
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}

	DB = pool
	fmt.Println("✅ Successfully connected to Neon Postgres")
}

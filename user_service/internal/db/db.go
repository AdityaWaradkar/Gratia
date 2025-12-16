package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"user_service/internal/config"
)

func Init(cfg *config.Config) *sql.DB {
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("DB ping failed:", err)
	}

	log.Println("Database connected")
	return db
}

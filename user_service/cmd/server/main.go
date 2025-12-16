package main

import (
	"user_service/internal/config"
	"user_service/internal/db"
	"user_service/internal/logger"
	"user_service/internal/server"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger.Init(cfg.LogLevel)

	// Initialize database
	dbConn := db.Init(cfg)

	// Start HTTP server
	server.Start(cfg, dbConn)
}

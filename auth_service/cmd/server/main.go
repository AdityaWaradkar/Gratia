package main

import (
	"context"
	"net/http"

	"github.com/adityawaradkar/gratia/auth_service/internal/auth"
	"github.com/adityawaradkar/gratia/auth_service/internal/config"
	"github.com/adityawaradkar/gratia/auth_service/internal/db"
	"github.com/adityawaradkar/gratia/auth_service/internal/logger"
	"github.com/adityawaradkar/gratia/auth_service/internal/server"
)

func main() {
	// Load configuration from environment variables
	config.Load()

	// Initialize logger with service name
	log := logger.New("AUTH")

	// Connect to the database
	dbPool := db.Connect(context.Background(), config.AppConfig.DatabaseURL)
	defer dbPool.Close()

	// Initialize repository for database operations
	repo := auth.NewRepository(dbPool)

	// Initialize service with business logic
	service := auth.NewService(
		repo,
		config.AppConfig.JWTSecret,
		config.AppConfig.AccessTokenTTL,
		config.AppConfig.RefreshTokenTTL,
	)

	// Initialize handler for HTTP requests
	handler := auth.NewHandler(service)

	// Register HTTP routes with the router
	router := server.RegisterRoutes(handler)

	// Start the HTTP server
	log.Printf("auth service running on port %s", config.AppConfig.Port)
	if err := http.ListenAndServe(":"+config.AppConfig.Port, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
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
	// Load configuration
	config.Load()

	// Initialize logger
	log := logger.New("AUTH")

	// Connect to database
	dbPool := db.Connect(context.Background(), config.AppConfig.DatabaseURL)
	defer dbPool.Close()

	// Initialize repository
	repo := auth.NewRepository(dbPool)

	// Initialize service
	service := auth.NewService(
		repo,
		config.AppConfig.JWTSecret,
		config.AppConfig.AccessTokenTTL,
		config.AppConfig.RefreshTokenTTL,
	)

	// Initialize handler
	handler := auth.NewHandler(service)

	// Register HTTP routes
	router := server.RegisterRoutes(handler)

	// Start HTTP server
	log.Printf("auth service running on port %s", config.AppConfig.Port)
	if err := http.ListenAndServe(":"+config.AppConfig.Port, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

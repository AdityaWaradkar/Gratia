package main

import (
	"context"
	"log"
	"net/http"

	"github.com/adityawaradkar/gratia/auth_service/internal/auth"
	"github.com/adityawaradkar/gratia/auth_service/internal/config"
	httpServer "github.com/adityawaradkar/gratia/auth_service/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to Postgres
	dbpool, err := pgxpool.New(context.Background(), config.GetDatabaseURL())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbpool.Close()

	// Initialize repository and service
	repo := auth.NewRepository(dbpool)
	service := auth.NewService(
		repo,
		config.GetJWTSecret(),
		config.GetAccessTokenTTL(),
		config.GetRefreshTokenTTL(),
	)

	// Initialize HTTP handler
	handler := auth.NewHandler(service)

	// Create the router with all routes
	router := httpServer.RegisterRoutes(handler)

	// Start server
	addr := ":" + config.GetPort()
	log.Printf("Auth service running on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

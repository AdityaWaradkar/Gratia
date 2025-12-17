package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/adityawaradkar/gratia/auth_service/internal/auth"
	"github.com/adityawaradkar/gratia/auth_service/internal/config"
	httpServer "github.com/adityawaradkar/gratia/auth_service/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Logger
	logger := log.New(os.Stdout, "[AUTH] ", log.LstdFlags|log.Lshortfile)

	// Database
	dbPool, err := pgxpool.New(context.Background(), config.GetDatabaseURL())
	if err != nil {
		logger.Fatalf("database connection failed: %v", err)
	}
	defer dbPool.Close()

	// Repository
	repo := auth.NewRepository(dbPool)

	// Service
	service := auth.NewService(
		repo,
		config.GetJWTSecret(),
		config.GetAccessTokenTTL(),
		config.GetRefreshTokenTTL(),
		config.AppConfig.UserServiceURL,
		logger,
	)

	// Handler
	handler := auth.NewHandler(service, dbPool)

	// Router
	router := httpServer.RegisterRoutes(handler)

	// Server
	address := ":" + config.GetPort()
	logger.Printf("auth service running on %s", address)

	if err := http.ListenAndServe(address, router); err != nil {
		logger.Fatalf("server stopped: %v", err)
	}
}

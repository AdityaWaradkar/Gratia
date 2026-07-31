package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adityawaradkar/gratia/claim_service/internal/claim"
	"github.com/adityawaradkar/gratia/claim_service/internal/client"
	"github.com/adityawaradkar/gratia/claim_service/internal/config"
	"github.com/adityawaradkar/gratia/claim_service/internal/db"
	"github.com/adityawaradkar/gratia/claim_service/internal/logger"
	authmw "github.com/adityawaradkar/gratia/claim_service/internal/middleware"
	"github.com/adityawaradkar/gratia/claim_service/internal/server"
)

func main() {
	// Load configuration from environment variables
	config.Load()
	cfg := config.AppConfig

	// Initialize logger with configured log level
	logr := logger.New(logger.Config{
		Level: cfg.LogLevel,
	})

	logr.Info(
		"starting claim service",
		"env", cfg.Env,
		"port", cfg.Port,
	)

	// Connect to the database
	dbConn, err := db.New()
	if err != nil {
		logr.Error(
			"failed to connect to database",
			"error", err,
		)
		os.Exit(1)
	}
	defer dbConn.Close()

	// Initialize repository for claim database operations
	repo := claim.NewClaimRepository(dbConn)

	// Initialize external service clients
	userClient := client.NewUserClient(
		cfg.UserServiceURL,
	)
	foodClient := client.NewFoodClient(
		cfg.FoodServiceURL,
	)

	// Initialize service layer with business logic
	claimService := claim.NewService(
		dbConn,
		repo,
		userClient,
		foodClient,
	)

	// Initialize handler for HTTP requests
	claimHandler := claim.NewHandler(
		claimService,
	)

	// Initialize authentication middleware
	authMiddleware := authmw.NewMiddleware(
		cfg.JWTSecret,
	)

	// Initialize router with all routes and middleware
	srv := server.NewServer(
		claimHandler,
		authMiddleware,
	)

	// Configure the HTTP server
	httpServer := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Start the HTTP server in a goroutine
	go func() {
		logr.Info(
			"http server started",
			"port", cfg.Port,
		)

		if err := httpServer.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			logr.Error(
				"http server failed",
				"error", err,
			)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	<-quit

	logr.Info("shutdown signal received")

	// Gracefully shutdown the server with timeout
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logr.Error(
			"server shutdown failed",
			"error", err,
		)
		os.Exit(1)
	}

	logr.Info("server exited cleanly")
}
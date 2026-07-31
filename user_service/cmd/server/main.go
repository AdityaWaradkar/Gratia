package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adityawaradkar/gratia/user_service/internal/config"
	"github.com/adityawaradkar/gratia/user_service/internal/db"
	"github.com/adityawaradkar/gratia/user_service/internal/logger"
	"github.com/adityawaradkar/gratia/user_service/internal/server"
	"github.com/adityawaradkar/gratia/user_service/internal/user"
)

func main() {
	// Load configuration from environment variables
	config.Load()

	// Initialize logger
	logger.Init()

	// Connect to the database
	database := db.Connect(config.AppConfig.DatabaseURL)
	defer database.Close()

	// Initialize repository, service, and handler
	repo := user.NewRepository(database)
	service := user.NewService(repo)
	handler := user.NewHandler(service)

	// Register HTTP routes
	router := server.RegisterRoutes(handler)

	// Configure the HTTP server
	srv := &http.Server{
		Addr:              ":" + config.AppConfig.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Start the HTTP server in a goroutine
	go func() {
		logger.Logger.Println(
			"user service running on port",
			config.AppConfig.Port,
		)

		if err := srv.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			logger.Logger.Fatalf("server failed: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)
	<-stop

	logger.Logger.Println("shutting down user service...")

	// Gracefully shutdown the server with timeout
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Logger.Fatalf("graceful shutdown failed: %v", err)
	}

	logger.Logger.Println("user service stopped")
}
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adityawaradkar/gratia/claim_service/internal/claim"
	"github.com/adityawaradkar/gratia/claim_service/internal/config"
	"github.com/adityawaradkar/gratia/claim_service/internal/db"
	"github.com/adityawaradkar/gratia/claim_service/internal/logger"
	"github.com/adityawaradkar/gratia/claim_service/internal/server"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize structured JSON logger
	log := logger.New("claim_service", cfg.LogLevel)

	// Establish context for startup tasks
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Info("starting claim service", slog.String("env", cfg.Env), slog.String("port", cfg.Port))

	// Connect to the database
	database, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer database.Close()

	// Initialize repository for claim database operations
	repo := claim.NewRepository(database)

	// Initialize external service clients
	userClient := claim.NewUserClient(cfg.UserServiceURL)
	foodClient := claim.NewFoodClient(cfg.FoodServiceURL)

	// Initialize service layer with business logic
	claimService := claim.NewService(repo, userClient, foodClient)

	// Initialize handler for HTTP requests
	claimHandler := claim.NewHandler(claimService)

	// Register routes with JWT authentication
	router := server.RegisterRoutes(claimHandler, cfg.JWTSecret)

	// Apply CORS middleware
	handlerWithCORS := corsMiddleware(router)

	// Configure the HTTP server
	httpServer := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handlerWithCORS,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Start the HTTP server in a goroutine
	go func() {
		log.Info("http server started", slog.String("port", cfg.Port))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutdown signal received")

	// Gracefully shutdown the server with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Info("server exited cleanly")
}

// corsMiddleware handles CORS headers for all responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow specific frontend origin
		allowedOrigin := "https://gratia-eta.vercel.app"

		// Allow localhost for development
		if r.Header.Get("Origin") == "http://localhost:3000" {
			allowedOrigin = "http://localhost:3000"
		}

		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
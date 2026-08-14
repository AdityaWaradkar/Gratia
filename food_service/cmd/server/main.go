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

	"github.com/adityawaradkar/gratia/food_service/internal/config"
	"github.com/adityawaradkar/gratia/food_service/internal/db"
	"github.com/adityawaradkar/gratia/food_service/internal/food"
	"github.com/adityawaradkar/gratia/food_service/internal/logger"
	"github.com/adityawaradkar/gratia/food_service/internal/server"
)

func main() {
	// Load configuration via dependency-friendly loading function
	cfg := config.Load()

	// Initialize structured JSON logger
	log := logger.New("food_service", cfg.LogLevel)

	// Establish context for startup tasks
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to the PostgreSQL connection pool securely
	database, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("failed to establish database connection", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer database.Close()

	// Initialize clean architectural layers
	repo := food.NewRepository(database)
	userClient := food.NewUserClient(cfg.UserServiceURL)
	service := food.NewService(repo, userClient)

	// Create a cancellable context to manage background routines lifecycle
	workerCtx, cancelWorker := context.WithCancel(context.Background())
	defer cancelWorker()

	// Start background worker for expiring listings asynchronously
	go startExpiryWorker(workerCtx, service, log)

	// Initialize HTTP handlers
	handler := food.NewHandler(service)

	// Register routes, injecting the explicitly loaded JWT secret
	router := server.RegisterRoutes(handler, cfg.JWTSecret)

	// Apply CORS middleware
	handlerWithCORS := corsMiddleware(router)

	// Configure robust server parameters to mitigate slow-client attacks
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handlerWithCORS,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Start the HTTP server in a goroutine
	go func() {
		log.Info("food service started successfully", slog.String("port", cfg.Port), slog.String("env", cfg.Env))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server failed unexpectedly", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Listen for OS interrupt signals for graceful shutdown execution
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Info("initiating graceful shutdown sequence...")

	// 1. Signal background workers to halt processing immediately
	cancelWorker()

	// 2. Allow active HTTP requests up to 10 seconds to finish before forcing termination
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("forced server shutdown encountered an error", slog.String("error", err.Error()))
	}

	log.Info("food service stopped cleanly")
}

// startExpiryWorker runs a background worker to expire food listings periodically
func startExpiryWorker(ctx context.Context, service *food.Service, log *slog.Logger) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("expiry worker shutting down")
			return

		case <-ticker.C:
			expireCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			if err := service.ExpireListings(expireCtx); err != nil {
				log.Error("expiry worker error", slog.String("error", err.Error()))
			}
			cancel()
		}
	}
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
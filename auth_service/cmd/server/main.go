package main

import (
    "context"
    "errors"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/adityawaradkar/gratia/auth_service/internal/auth"
    "github.com/adityawaradkar/gratia/auth_service/internal/config"
    "github.com/adityawaradkar/gratia/auth_service/internal/db"
    "github.com/adityawaradkar/gratia/auth_service/internal/logger"
    "github.com/adityawaradkar/gratia/auth_service/internal/server"
)

func main() {
    // Load configuration securely without relying on global variables
    cfg := config.Load()

    // Initialize structured JSON logging based on the configured environment level
    log := logger.New("auth_service", cfg.LogLevel)

    // Create a base context for application startup operations
    ctx := context.Background()

    // Establish a highly concurrent connection pool to PostgreSQL and handle errors safely
    dbPool, err := db.Connect(ctx, cfg.DatabaseURL)
    if err != nil {
        log.Error("Failed to initialize database pool", "error", err)
        os.Exit(1)
    }
    defer dbPool.Close()

    // Wire up the repository data access layer with the active connection pool
    repo := auth.NewRepository(dbPool)

    // Inject required dependencies into the core authentication business logic
    service := auth.NewService(repo, cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)

    // Map the core service layer to the HTTP transport handlers
    handler := auth.NewHandler(service)

    // Construct the strict routing tree and inject the JWT secret for middleware verification
    router := server.RegisterRoutes(handler, cfg.JWTSecret)

    // Apply CORS middleware
    handlerWithCORS := corsMiddleware(router)

    // Configure the HTTP server with explicit timeouts to prevent resource exhaustion
    srv := &http.Server{
        Addr:         ":" + cfg.Port,
        Handler:      handlerWithCORS,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Run the server in a separate goroutine to allow for non-blocking signal catching
    go func() {
        log.Info("Starting server", "port", cfg.Port, "env", cfg.Env)
        if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
            log.Error("Server encountered a fatal error", "error", err)
            os.Exit(1)
        }
    }()

    // Listen for standard termination signals from container orchestrators like Docker
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Info("Shutting down server gracefully...")

    // Provide a strict 10-second window for active network requests to finish completely
    shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    // Execute the graceful shutdown and catch any dangling connection errors
    if err := srv.Shutdown(shutdownCtx); err != nil {
        log.Error("Server forced to shutdown abruptly", "error", err)
        os.Exit(1)
    }

    log.Info("Server exited properly")
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
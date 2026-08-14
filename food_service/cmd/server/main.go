package main

import (
    "context"
    "errors"
    "log"
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

    // Initialize application logger
    logger.Init()

    // Connect to the database connection pool
    pool := db.Connect(cfg.DatabaseURL)
    defer pool.Close()

    // Initialize clean architectural layers
    repo := food.NewRepository(pool)
    userClient := food.NewUserClient(cfg.UserServiceURL)
    service := food.NewService(repo, userClient)

    // Create a cancellable context to manage background routines lifecycle
    workerCtx, cancelWorker := context.WithCancel(context.Background())
    defer cancelWorker()

    // Start background worker for expiring listings asynchronously
    go startExpiryWorker(workerCtx, service)

    // Initialize HTTP handlers
    handler := food.NewHandler(service)

    // Register routes, injecting the explicitly loaded JWT secret
    router := server.RegisterRoutes(handler, cfg.JWTSecret)

    // Configure robust server parameters to mitigate slow-client attacks
    srv := &http.Server{
        Addr:              ":" + cfg.Port,
        Handler:           router,
        ReadHeaderTimeout: 5 * time.Second,
        ReadTimeout:       15 * time.Second,
        WriteTimeout:      15 * time.Second,
        IdleTimeout:       60 * time.Second,
    }

    // Start the HTTP server in a goroutine
    go func() {
        log.Printf("food_service started successfully on port %s", cfg.Port)
        if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
            log.Fatalf("server failed unexpectedly: %v", err)
        }
    }()

    // Listen for OS interrupt signals for graceful shutdown execution
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
    <-stop

    log.Println("initiating graceful shutdown sequence...")

    // 1. Signal background workers to halt processing immediately
    cancelWorker()

    // 2. Allow active HTTP requests up to 10 seconds to finish before forcing termination
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer shutdownCancel()

    if err := srv.Shutdown(shutdownCtx); err != nil {
        log.Printf("forced server shutdown encountered an error: %v", err)
    }

    log.Println("food_service stopped cleanly")
}

// startExpiryWorker runs a background worker to expire food listings periodically
func startExpiryWorker(ctx context.Context, service *food.Service) {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done(): // Graceful termination signal received
            log.Println("expiry worker shutting down")
            return
            
        case <-ticker.C:
            // Use a bounded context for the actual database operation
            expireCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
            if err := service.ExpireListings(expireCtx); err != nil {
                log.Printf("expiry worker error: %v", err)
            }
            cancel() // Free resources immediately after the check finishes
        }
    }
}
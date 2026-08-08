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

    "github.com/adityawaradkar/gratia/user_service/internal/config"
    "github.com/adityawaradkar/gratia/user_service/internal/db"
	"github.com/adityawaradkar/gratia/user_service/internal/logger"
    "github.com/adityawaradkar/gratia/user_service/internal/server"
    "github.com/adityawaradkar/gratia/user_service/internal/user"
)

func main() {
    // Load configuration via dependency-friendly loading function
    cfg := config.Load()

    // Initialize structured JSON logger
    log := logger.New("user_service", cfg.LogLevel)

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

    // Initialize clean architectural layers (Repository -> Service -> Handler)
    repo := user.NewRepository(database)
    service := user.NewService(repo)
    handler := user.NewHandler(service)

    // Register HTTP routes, passing the required JWT secret for middleware verification
    router := server.RegisterRoutes(handler, cfg.JWTSecret)

    // Configure the production HTTP server
    srv := &http.Server{
        Addr:              ":" + cfg.Port,
        Handler:           router,
        ReadHeaderTimeout: 5 * time.Second,
        ReadTimeout:       15 * time.Second,
        WriteTimeout:      15 * time.Second,
        IdleTimeout:       60 * time.Second,
    }

    // Start the HTTP server asynchronously
    go func() {
        log.Info("user service started successfully", slog.String("port", cfg.Port), slog.String("env", cfg.Env))
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

    // Allow active requests 10 seconds to finish processing before forcing termination
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer shutdownCancel()

    if err := srv.Shutdown(shutdownCtx); err != nil {
        log.Error("forced server shutdown encountered an error", slog.String("error", err.Error()))
    }

    log.Info("user service stopped cleanly")
}
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	"github.com/adityawaradkar/gratia-auth/internal/config"
	handler "github.com/adityawaradkar/gratia-auth/internal/handlers"
	"github.com/adityawaradkar/gratia-auth/internal/middleware"
	"github.com/adityawaradkar/gratia-auth/internal/repository"
	"github.com/adityawaradkar/gratia-auth/internal/service"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found or failed to load")
	}

	cfg := config.LoadConfig()

	// Connect to PostgreSQL database
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authSvc)

	// Create rate limiter: 5 requests per minute burst 10, data expires after 5 minutes
	rateLimiter := middleware.NewRateLimiter(5, 10, 5*time.Minute)

	mux := http.NewServeMux()

	// Protect /register and /login with rate limiter middleware
	mux.Handle("/register", rateLimiter.Limit(http.HandlerFunc(authHandler.Register)))
	mux.Handle("/login", rateLimiter.Limit(http.HandlerFunc(authHandler.Login)))

	// Other auth routes without rate limiting
	mux.Handle("/api/v1/auth/", http.StripPrefix("/api/v1/auth", authHandler.Routes()))

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      withCORS(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		fmt.Printf("auth-service is running on port %s\n", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	fmt.Println("\nShutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	fmt.Println("Server gracefully stopped.")
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/adityawaradkar/gratia/food_service/internal/config"
	"github.com/adityawaradkar/gratia/food_service/internal/db"
	"github.com/adityawaradkar/gratia/food_service/internal/food"
	"github.com/adityawaradkar/gratia/food_service/internal/logger"
	"github.com/adityawaradkar/gratia/food_service/internal/server"
)

func main() {
	// Load configuration
	config.Load()

	// Initialize logger
	logger.Init()

	// Connect to database
	pool := db.Connect(config.AppConfig.DatabaseURL)
	defer pool.Close()

	// Initialize repository
	repo := food.NewRepository(pool)

	// Initialize user service client
	userClient := food.NewUserClient(config.AppConfig.UserServiceURL)

	// Initialize service
	service := food.NewService(repo, userClient)

	// Start expiry worker
	go startExpiryWorker(service)

	// Initialize handler
	handler := food.NewHandler(service)

	// Register routes
	router := server.RegisterRoutes(handler)

	// Start HTTP server
	log.Printf("food_service running on port %s", config.AppConfig.Port)

	if err := http.ListenAndServe(
		":"+config.AppConfig.Port,
		router,
	); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func startExpiryWorker(service *food.Service) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := service.ExpireListings(
			context.Background(),
		); err != nil {
			log.Printf(
				"expiry worker error: %v",
				err,
			)
		}
	}
}
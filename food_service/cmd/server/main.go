package main

import (
	"log"
	"net/http"

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

	// Initialize handler
	handler := food.NewHandler(service)

	// Register routes
	router := server.RegisterRoutes(handler)

	// Start HTTP server
	log.Printf("food_service running on port %s", config.AppConfig.Port)
	if err := http.ListenAndServe(":"+config.AppConfig.Port, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

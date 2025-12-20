package main

import (
	"net/http"

	"github.com/adityawaradkar/gratia/user_service/internal/config"
	"github.com/adityawaradkar/gratia/user_service/internal/db"
	"github.com/adityawaradkar/gratia/user_service/internal/logger"
	"github.com/adityawaradkar/gratia/user_service/internal/server"
	"github.com/adityawaradkar/gratia/user_service/internal/user"
)

func main() {
	// Load configuration
	config.Load()

	// Initialize logger
	logger.Init()

	// Connect to database
	database := db.Connect(config.AppConfig.DatabaseURL)
	defer database.Close()

	// Wire domain
	repo := user.NewRepository(database)
	service := user.NewService(repo)
	handler := user.NewHandler(service)

	// Setup HTTP server
	router := server.RegisterRoutes(handler)

	logger.Logger.Println("user service running on port", config.AppConfig.Port)

	// Start server
	if err := http.ListenAndServe(":"+config.AppConfig.Port, router); err != nil {
		logger.Logger.Fatalf("server failed: %v", err)
	}
}

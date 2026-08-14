package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Env            string
	Port           string
	JWTSecret      string
	DatabaseURL    string
	UserServiceURL string
	FoodServiceURL string
	LogLevel       string
}

// Load loads application configuration from environment variables
func Load() *Config {
	_ = godotenv.Load()

	jwtSecret := mustEnv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		log.Fatal("Fatal: JWT_SECRET must be at least 32 characters")
	}

	return &Config{
		Env:            getEnv("ENV", "development"),
		Port:           getEnv("PORT", "8083"),
		JWTSecret:      jwtSecret,
		DatabaseURL:    mustEnv("DATABASE_URL"),
		UserServiceURL: mustEnv("USER_SERVICE_URL"),
		FoodServiceURL: mustEnv("FOOD_SERVICE_URL"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv retrieves an environment variable with a default value
func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// mustEnv retrieves a required environment variable or exits with an error
func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Fatal: %s environment variable is strictly required", key)
	}
	return value
}
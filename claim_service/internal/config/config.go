package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Env string

	Port string

	JWTSecret string

	DatabaseURL string

	UserServiceURL string
	FoodServiceURL string

	LogLevel string
}

// AppConfig is the global configuration instance
var AppConfig *Config

// Load loads application configuration from environment variables
func Load() {
	_ = godotenv.Load()

	AppConfig = &Config{
		Env: getEnv(
			"ENV",
			"development",
		),

		Port: getEnv(
			"PORT",
			"8083",
		),

		JWTSecret: mustEnv(
			"JWT_SECRET",
		),

		DatabaseURL: mustEnv(
			"DATABASE_URL",
		),

		UserServiceURL: mustEnv(
			"USER_SERVICE_URL",
		),

		FoodServiceURL: mustEnv(
			"FOOD_SERVICE_URL",
		),

		LogLevel: getEnv(
			"LOG_LEVEL",
			"info",
		),
	}
}

// getEnv retrieves an environment variable with a default value
func getEnv(
	key string,
	defaultValue string,
) string {

	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

// mustEnv retrieves a required environment variable or exits with an error
func mustEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("%s is required", key)
	}

	return value
}
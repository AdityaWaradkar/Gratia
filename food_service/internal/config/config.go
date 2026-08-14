package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds food service configuration
type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	UserServiceURL string
	Env            string
	LogLevel       string
}

// Load reads environment variables and returns a Config instance
func Load() *Config {
	_ = godotenv.Load()

	jwtSecret := mustEnv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		log.Fatal("Fatal: JWT_SECRET must be at least 32 characters")
	}

	return &Config{
		Port:           getEnv("PORT", "8082"),
		DatabaseURL:    mustEnv("DATABASE_URL"),
		JWTSecret:      jwtSecret,
		UserServiceURL: getEnv("USER_SERVICE_URL", "http://user_service:8081"),
		Env:            getEnv("ENV", "development"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv retrieves an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// mustEnv retrieves a required environment variable or exits with an error
func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Fatal: %s environment variable is strictly required", key)
	}
	return v
}
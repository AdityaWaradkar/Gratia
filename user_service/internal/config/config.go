package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds user service configuration
type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

// AppConfig is the loaded configuration
var AppConfig *Config

// Load reads environment variables into Config
func Load() {
	_ = godotenv.Load()

	AppConfig = &Config{
		Port:        getEnv("PORT", "8081"),
		DatabaseURL: mustEnv("DATABASE_URL"),
		JWTSecret:   mustEnv("JWT_SECRET"),
	}
}

// getEnv reads optional env variable
func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// mustEnv reads required env variable
func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("%s is required", key)
	}
	return v
}

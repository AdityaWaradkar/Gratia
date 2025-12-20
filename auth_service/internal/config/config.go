package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds service configuration
type Config struct {
	Port             string
	DatabaseURL      string
	JWTSecret        string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
}

// AppConfig is the loaded configuration
var AppConfig *Config

// Load loads environment variables into Config
func Load() {
	_ = godotenv.Load()

	AppConfig = &Config{
		Port:            getEnv("PORT", "8080"),
		DatabaseURL:     mustEnv("DATABASE_URL"),
		JWTSecret:       mustEnv("JWT_SECRET"),
		AccessTokenTTL:  time.Duration(getIntEnv("ACCESS_TOKEN_MINUTES", 15)) * time.Minute,
		RefreshTokenTTL: time.Duration(getIntEnv("REFRESH_TOKEN_DAYS", 30)) * 24 * time.Hour,
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

// getIntEnv reads int env variable
func getIntEnv(key string, defaultValue int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("invalid value for %s", key)
	}
	return i
}

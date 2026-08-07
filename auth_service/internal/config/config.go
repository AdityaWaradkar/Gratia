package config

import (
    "log"
    "os"
    "strconv"
    "time"

    "github.com/joho/godotenv"
)

// Config holds all the environment-specific settings for the service
type Config struct {
    Port            string
    DatabaseURL     string
    JWTSecret       string
    AccessTokenTTL  time.Duration
    RefreshTokenTTL time.Duration
    Env             string // Added to support environment tracking (e.g., development or production)
    LogLevel        string // Added to configure structured logging output
}

// Load reads environment variables and returns a Config instance instead of mutating a global state
func Load() *Config {
    _ = godotenv.Load() // Errors are ignored because production environments often inject variables directly without a .env file

    return &Config{
        Port:            getEnv("PORT", "8080"),
        DatabaseURL:     mustEnv("DATABASE_URL"),
        JWTSecret:       mustEnv("JWT_SECRET"),
        AccessTokenTTL:  time.Duration(getIntEnv("ACCESS_TOKEN_MINUTES", 15)) * time.Minute,
        RefreshTokenTTL: time.Duration(getIntEnv("REFRESH_TOKEN_DAYS", 30)) * 24 * time.Hour,
        Env:             getEnv("ENV", "development"),
        LogLevel:        getEnv("LOG_LEVEL", "info"),
    }
}

// getEnv reads an optional environment variable or returns the provided default value
func getEnv(key, defaultValue string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return defaultValue
}

// mustEnv reads a strictly required environment variable and halts the application if it is missing
func mustEnv(key string) string {
    v := os.Getenv(key)
    if v == "" {
        log.Fatalf("Fatal: %s environment variable is strictly required", key)
    }
    return v
}

// getIntEnv reads an integer environment variable or returns the provided default value
func getIntEnv(key string, defaultValue int) int {
    v := os.Getenv(key)
    if v == "" {
        return defaultValue
    }
    i, err := strconv.Atoi(v)
    if err != nil {
        log.Fatalf("Fatal: invalid integer value provided for %s", key)
    }
    return i
}
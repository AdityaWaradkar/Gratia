package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

// Config holds all the environment-specific settings for the user service
type Config struct {
    Port        string
    DatabaseURL string
    JWTSecret   string
    Env         string
    LogLevel    string
}

// Load reads environment variables and returns a Config instance instead of mutating a global state
func Load() *Config {
    _ = godotenv.Load()

    return &Config{
        Port:        getEnv("PORT", "8081"),
        DatabaseURL: mustEnv("DATABASE_URL"),
        JWTSecret:   mustEnv("JWT_SECRET"),
        Env:         getEnv("ENV", "development"),
        LogLevel:    getEnv("LOG_LEVEL", "info"),
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
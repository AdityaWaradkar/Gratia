package config

import (
    "log"
    "os"
    "strconv"
    "time"

    "github.com/joho/godotenv"
)

// Config holds all the environment-specific settings for the user service
type Config struct {
    Port            string
    DatabaseURL     string
    JWTSecret       string
    Env             string
    LogLevel        string
    AuthServiceURL  string
}

// Load reads environment variables and returns a Config instance
func Load() *Config {
    _ = godotenv.Load()

    jwtSecret := mustEnv("JWT_SECRET")
    if len(jwtSecret) < 32 {
        log.Fatal("Fatal: JWT_SECRET must be at least 32 characters")
    }

    return &Config{
        Port:           getEnv("PORT", "8081"),
        DatabaseURL:    mustEnv("DATABASE_URL"),
        JWTSecret:      jwtSecret,
        Env:            getEnv("ENV", "development"),
        LogLevel:       getEnv("LOG_LEVEL", "info"),
        AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://auth_service:8080"),
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
        log.Fatalf("Fatal: invalid integer value provided for %s: %s", key, v)
    }
    return i
}

// getDurationEnv reads a duration environment variable or returns the provided default value
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
    v := os.Getenv(key)
    if v == "" {
        return defaultValue
    }
    d, err := time.ParseDuration(v)
    if err != nil {
        log.Fatalf("Fatal: invalid duration value provided for %s: %s", key, v)
    }
    return d
}
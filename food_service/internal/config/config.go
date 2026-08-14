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
}

// AppConfig is the loaded configuration
var AppConfig *Config

// Load reads environment variables into Config
func Load() {
    _ = godotenv.Load()

    AppConfig = &Config{
        Port:           getEnv("PORT", "8082"),
        DatabaseURL:    mustEnv("DATABASE_URL"),
        JWTSecret:      mustEnv("JWT_SECRET"),
        UserServiceURL: getEnv("USER_SERVICE_URL", "http://localhost:8081"),
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
        log.Fatalf("%s is required", key)
    }
    return v
}
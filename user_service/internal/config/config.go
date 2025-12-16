package config

import (
	"log"
	"os"
)

type Config struct {
	ServerPort string
	DBUrl      string
	LogLevel   string
	JWTSecret  string
}

func Load() *Config {
	cfg := &Config{
		ServerPort: getEnv("SERVER_PORT", "8081"),
		DBUrl:      getEnv("DB_URL", ""),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
		JWTSecret:  getEnv("JWT_SECRET", ""),
	}

	if cfg.DBUrl == "" {
		log.Fatal("DB_URL is required")
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

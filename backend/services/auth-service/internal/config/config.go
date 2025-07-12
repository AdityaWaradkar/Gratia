package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string
	JWTSecret    string
	SupabaseURL  string
	SupabaseAnon string
	ServerPort   string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found in auth-service folder, reading from environment variables")
	}

	config := &Config{
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		ServerPort:   os.Getenv("SERVER_PORT"),
	}

	if config.ServerPort == "" {
		config.ServerPort = "8081"
	}

	if config.DatabaseURL == "" || config.JWTSecret == "" {
		log.Fatal("DATABASE_URL and JWT_SECRET must be set")
	}

	return config
}

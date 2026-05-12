package config

import (
	"log"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	AppPort string `validate:"required"`

	DatabaseURL string `validate:"required,url"`

	JWTSecret string `validate:"required,min=32"`

	JWTExpirationMinutes int `validate:"required,gte=15,lte=1440"`

	AdminEmail string `validate:"required,email"`

	AdminPassword string `validate:"required,min=8"`
}

func LoadConfig() *Config {
	LoadEnv()

	expirationMinutes, err := strconv.Atoi(
		GetEnv("JWT_EXPIRATION_MINUTES"),
	)

	if err != nil {
		log.Fatal("Invalid JWT_EXPIRATION_MINUTES")
	}

	cfg := &Config{
		AppPort: GetEnv("APP_PORT"),

		DatabaseURL: GetEnv("DATABASE_URL"),

		JWTSecret: GetEnv("JWT_SECRET"),

		JWTExpirationMinutes: expirationMinutes,

		AdminEmail: GetEnv("ADMIN_EMAIL"),

		AdminPassword: GetEnv("ADMIN_PASSWORD"),
	}

	validate := validator.New()

	err = validate.Struct(cfg)

	if err != nil {
		log.Fatalf(
			"Invalid configuration: %v",
			err,
		)
	}

	return cfg
}
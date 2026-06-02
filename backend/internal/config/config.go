package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppName string
	AppEnv  string
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret           string
	JWTExpirationMinutes int
}

func Load() (Config, error) {
	cfg := Config{
		AppName: envOrDefault("APP_NAME", "law-office-api"),
		AppEnv:  envOrDefault("APP_ENV", "local"),
		AppPort: envOrDefault("APP_PORT", "8080"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     envOrDefault("DB_PORT", "5432"),
		DBUser:     envOrDefault("DB_USER", "postgres"),
		DBPassword: envOrDefault("DB_PASSWORD", "postgres"),
		DBName:     envOrDefault("DB_NAME", "legalflow"),

		JWTSecret:           os.Getenv("JWT_SECRET"),
		JWTExpirationMinutes: envOrDefaultInt("JWT_EXPIRATION_MINUTES", 60),
	}

	if cfg.DBHost == "" {
		return Config{}, fmt.Errorf("required environment variable DB_HOST is not set")
	}

	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("required environment variable JWT_SECRET is not set")
	}

	return cfg, nil
}

func envOrDefaultInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		var result int
		fmt.Sscanf(v, "%d", &result)
		return result
	}
	return fallback
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	cfg := &Config{
		HTTPAddr:    getEnv("HTTP_ADDR", ":1323"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

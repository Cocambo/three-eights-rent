package config

import (
	"fmt"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort      string
	GinMode         string
	JWTAccessSecret string
	UserServiceURL  string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		GinMode:         getEnv("GIN_MODE", "debug"),
		JWTAccessSecret: os.Getenv("JWT_ACCESS_SECRET"),
		UserServiceURL:  os.Getenv("USER_SERVICE_URL"),
	}

	if cfg.JWTAccessSecret == "" {
		return nil, fmt.Errorf("JWT_ACCESS_SECRET is required")
	}

	if cfg.UserServiceURL == "" {
		return nil, fmt.Errorf("USER_SERVICE_URL is required")
	}

	parsedURL, err := url.Parse(cfg.UserServiceURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("USER_SERVICE_URL must be a valid absolute URL")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

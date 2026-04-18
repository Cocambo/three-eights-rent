package config

import (
	"car-service/database"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string
	Server ServerConfig
	DB     database.Config
	JWT    JWTConfig
	MinIO  MinIOConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type JWTConfig struct {
	AccessSecret string
}

type MinIOConfig struct {
	Endpoint      string
	AccessKey     string
	SecretKey     string
	UseSSL        bool
	DefaultBucket string
	PresignTTL    time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	dbConfig, err := loadDatabaseConfig()
	if err != nil {
		return nil, err
	}

	minioPresignTTL, err := getEnvAsDuration("MINIO_PRESIGN_TTL", 15*time.Minute)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		AppEnv: getEnv("APP_ENV", "development"),
		Server: ServerConfig{
			Port:    getEnv("SERVER_PORT", "8082"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		DB: dbConfig,
		JWT: JWTConfig{
			AccessSecret: getEnv("JWT_ACCESS_SECRET", "dev-access-secret"),
		},
		MinIO: MinIOConfig{
			Endpoint:      getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey:     getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey:     getEnv("MINIO_SECRET_KEY", "minioadmin"),
			UseSSL:        getEnvAsBool("MINIO_USE_SSL", false),
			DefaultBucket: getEnv("MINIO_DEFAULT_BUCKET", getEnv("MINIO_BUCKET", "car-images")),
			PresignTTL:    minioPresignTTL,
		},
	}

	if cfg.Server.Port == "" {
		return nil, fmt.Errorf("SERVER_PORT is required")
	}

	if err := cfg.DB.Validate(); err != nil {
		return nil, fmt.Errorf("database config: %w", err)
	}

	if cfg.MinIO.Endpoint == "" {
		return nil, fmt.Errorf("MINIO_ENDPOINT is required")
	}

	if cfg.MinIO.DefaultBucket == "" {
		return nil, fmt.Errorf("MINIO_DEFAULT_BUCKET is required")
	}

	if cfg.MinIO.PresignTTL <= 0 {
		return nil, fmt.Errorf("MINIO_PRESIGN_TTL must be greater than zero")
	}

	return cfg, nil
}

func loadDatabaseConfig() (database.Config, error) {
	maxIdleConns, err := getEnvAsInt("DB_MAX_IDLE_CONNS", 10)
	if err != nil {
		return database.Config{}, err
	}

	maxOpenConns, err := getEnvAsInt("DB_MAX_OPEN_CONNS", 25)
	if err != nil {
		return database.Config{}, err
	}

	connMaxLifetime, err := getEnvAsDuration("DB_CONN_MAX_LIFETIME", time.Hour)
	if err != nil {
		return database.Config{}, err
	}

	connMaxIdleTime, err := getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", 15*time.Minute)
	if err != nil {
		return database.Config{}, err
	}

	return database.Config{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnv("DB_PORT", "5432"),
		User:            getEnv("DB_USER", "postgres"),
		Password:        getEnv("DB_PASSWORD", "postgres"),
		Name:            getEnv("DB_NAME", "car_service"),
		SSLMode:         getEnv("DB_SSLMODE", "disable"),
		MaxIdleConns:    maxIdleConns,
		MaxOpenConns:    maxOpenConns,
		ConnMaxLifetime: connMaxLifetime,
		ConnMaxIdleTime: connMaxIdleTime,
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value == "true" || value == "1"
}

func getEnvAsInt(key string, fallback int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid integer: %w", key, err)
	}

	return parsed, nil
}

func getEnvAsDuration(key string, fallback time.Duration) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid duration: %w", key, err)
	}

	return parsed, nil
}

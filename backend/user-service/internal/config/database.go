package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg *Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cfg.PostgresDSN()), &gorm.Config{})
}

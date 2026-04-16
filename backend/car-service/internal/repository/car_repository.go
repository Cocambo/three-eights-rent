package repository

import (
	"context"

	"car-service/internal/domains"

	"gorm.io/gorm"
)

type CarRepository interface {
	List(ctx context.Context) ([]domains.Car, error)
}

type carRepository struct {
	db *gorm.DB
}

func NewCarRepository(db *gorm.DB) CarRepository {
	return &carRepository{db: db}
}

func (r *carRepository) List(ctx context.Context) ([]domains.Car, error) {
	var cars []domains.Car

	if err := r.db.WithContext(ctx).Order("id DESC").Find(&cars).Error; err != nil {
		return nil, err
	}

	return cars, nil
}

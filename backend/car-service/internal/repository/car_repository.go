package repository

import (
	"context"
	"errors"

	"car-service/internal/domains"
	apperrors "car-service/internal/errors"

	"gorm.io/gorm"
)

type CarRepository interface {
	List(ctx context.Context) ([]domains.Car, error)
	GetByID(ctx context.Context, id uint) (domains.Car, error)
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

func (r *carRepository) GetByID(ctx context.Context, id uint) (domains.Car, error) {
	var car domains.Car

	if err := r.db.WithContext(ctx).First(&car, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domains.Car{}, apperrors.New(apperrors.ErrNotFound, "car not found")
		}

		return domains.Car{}, err
	}

	return car, nil
}

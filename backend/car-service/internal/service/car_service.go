package service

import (
	"context"

	"car-service/internal/domains"
	apperrors "car-service/internal/errors"
	"car-service/internal/repository"
)

type CarService interface {
	List(ctx context.Context) ([]domains.Car, error)
	GetByID(ctx context.Context, id uint) (domains.Car, error)
}

type carService struct {
	carRepository repository.CarRepository
}

func NewCarService(carRepository repository.CarRepository) CarService {
	return &carService{carRepository: carRepository}
}

func (s *carService) List(ctx context.Context) ([]domains.Car, error) {
	return s.carRepository.List(ctx)
}

func (s *carService) GetByID(ctx context.Context, id uint) (domains.Car, error) {
	if id == 0 {
		return domains.Car{}, apperrors.New(apperrors.ErrValidation, "invalid car id")
	}

	return s.carRepository.GetByID(ctx, id)
}

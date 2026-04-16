package service

import (
	"context"

	"car-service/internal/domains"
	"car-service/internal/repository"
)

type CarService interface {
	List(ctx context.Context) ([]domains.Car, error)
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

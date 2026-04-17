package service

import (
	"context"

	"car-service/internal/domains"
	"car-service/internal/dto"
	apperrors "car-service/internal/errors"
	"car-service/internal/repository"
)

type CarService interface {
	List(ctx context.Context, query dto.ListCarsQuery) ([]domains.Car, int64, error)
	GetByID(ctx context.Context, id uint) (domains.Car, error)
}

type carService struct {
	carRepository repository.CarRepository
}

func NewCarService(carRepository repository.CarRepository) CarService {
	return &carService{carRepository: carRepository}
}

func (s *carService) List(ctx context.Context, query dto.ListCarsQuery) ([]domains.Car, int64, error) {
	filter := repository.CarFilter{
		Search:       query.Q,
		Brand:        query.Brand,
		Model:        query.Model,
		YearFrom:     query.YearFrom,
		YearTo:       query.YearTo,
		FuelType:     query.FuelType,
		Transmission: query.Transmission,
		BodyType:     query.BodyType,
		SeatsMin:     query.SeatsMin,
		PriceMin:     query.PriceMin,
		PriceMax:     query.PriceMax,
		Purpose:      query.Purpose,
		SortBy:       query.SortBy,
		SortOrder:    query.SortOrder,
		Limit:        query.LimitValue(),
		Offset:       query.OffsetValue(),
	}

	total, err := s.carRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	cars, err := s.carRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return cars, total, nil
}

func (s *carService) GetByID(ctx context.Context, id uint) (domains.Car, error) {
	if id == 0 {
		return domains.Car{}, apperrors.New(apperrors.ErrValidation, "invalid car id")
	}

	return s.carRepository.GetByID(ctx, id)
}

package service

import (
	"context"

	"car-service/internal/domains"
	apperrors "car-service/internal/errors"
	"car-service/internal/repository"
)

type FavoriteService interface {
	GetUserFavorites(ctx context.Context, userID uint) ([]domains.Favorite, error)
	AddToFavorites(ctx context.Context, userID, carID uint) error
	RemoveFromFavorites(ctx context.Context, userID, carID uint) error
}

type favoriteService struct {
	favoriteRepository repository.FavoriteRepository
	carRepository      repository.CarRepository
}

func NewFavoriteService(
	favoriteRepository repository.FavoriteRepository,
	carRepository repository.CarRepository,
) FavoriteService {
	return &favoriteService{
		favoriteRepository: favoriteRepository,
		carRepository:      carRepository,
	}
}

func (s *favoriteService) GetUserFavorites(ctx context.Context, userID uint) ([]domains.Favorite, error) {
	if userID == 0 {
		return nil, validationError("user_id must be greater than zero")
	}

	favorites, err := s.favoriteRepository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	for i := range favorites {
		favorites[i].Car.CarImages = catalogImages(favorites[i].Car.CarImages)
	}

	return favorites, nil
}

func (s *favoriteService) AddToFavorites(ctx context.Context, userID, carID uint) error {
	if userID == 0 {
		return validationError("user_id must be greater than zero")
	}

	if carID == 0 {
		return validationError("car_id must be greater than zero")
	}

	exists, err := s.carRepository.ExistsByID(ctx, carID)
	if err != nil {
		return err
	}

	if !exists {
		return apperrors.New(apperrors.ErrNotFound, "car not found")
	}

	return s.favoriteRepository.Add(ctx, userID, carID)
}

func (s *favoriteService) RemoveFromFavorites(ctx context.Context, userID, carID uint) error {
	if userID == 0 {
		return validationError("user_id must be greater than zero")
	}

	if carID == 0 {
		return validationError("car_id must be greater than zero")
	}

	return s.favoriteRepository.Remove(ctx, userID, carID)
}

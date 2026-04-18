package service

import (
	"context"
	"time"

	apperrors "car-service/internal/errors"
	"car-service/internal/repository"
	"car-service/internal/storage"
)

type FavoriteItem struct {
	CarID   uint
	AddedAt time.Time
	Car     CatalogCar
}

type FavoriteService interface {
	GetUserFavorites(ctx context.Context, userID uint) ([]FavoriteItem, error)
	AddToFavorites(ctx context.Context, userID, carID uint) error
	RemoveFromFavorites(ctx context.Context, userID, carID uint) error
}

type favoriteService struct {
	favoriteRepository repository.FavoriteRepository
	carRepository      repository.CarRepository
	imageStorage       storage.ImageStorageService
}

func NewFavoriteService(
	favoriteRepository repository.FavoriteRepository,
	carRepository repository.CarRepository,
	imageStorage storage.ImageStorageService,
) FavoriteService {
	return &favoriteService{
		favoriteRepository: favoriteRepository,
		carRepository:      carRepository,
		imageStorage:       imageStorage,
	}
}

func (s *favoriteService) GetUserFavorites(ctx context.Context, userID uint) ([]FavoriteItem, error) {
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

	items := make([]FavoriteItem, 0, len(favorites))
	for _, favorite := range favorites {
		car, err := toCatalogCar(ctx, s.imageStorage, favorite.Car)
		if err != nil {
			return nil, err
		}

		items = append(items, FavoriteItem{
			CarID:   favorite.CarID,
			AddedAt: favorite.CreatedAt,
			Car:     car,
		})
	}

	return items, nil
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

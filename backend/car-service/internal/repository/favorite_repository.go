package repository

import (
	"context"
	"errors"

	"car-service/internal/domains"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FavoriteRepository interface {
	Add(ctx context.Context, userID, carID uint) error
	Remove(ctx context.Context, userID, carID uint) error
	ListByUserID(ctx context.Context, userID uint) ([]domains.Favorite, error)
	Exists(ctx context.Context, userID, carID uint) (bool, error)
}

type gormFavoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) FavoriteRepository {
	return &gormFavoriteRepository{db: db}
}

func (r *gormFavoriteRepository) Add(ctx context.Context, userID, carID uint) error {
	favorite := domains.Favorite{
		UserID: userID,
		CarID:  carID,
	}

	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "user_id"},
				{Name: "car_id"},
			},
			DoNothing: true,
		}).
		Create(&favorite).Error
	if err != nil {
		return mapRepositoryError(err, "favorite car not found")
	}

	return nil
}

func (r *gormFavoriteRepository) Remove(ctx context.Context, userID, carID uint) error {
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND car_id = ?", userID, carID).
		Delete(&domains.Favorite{}).Error
	if err != nil {
		return mapRepositoryError(err, "favorite not found")
	}

	return nil
}

func (r *gormFavoriteRepository) ListByUserID(ctx context.Context, userID uint) ([]domains.Favorite, error) {
	var favorites []domains.Favorite

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Preload("Car").
		Preload("Car.CarImages", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("is_main = ?", true).
				Order("is_main DESC").
				Order("sort_order ASC").
				Order("id ASC")
		}).
		Find(&favorites).Error
	if err != nil {
		return nil, mapRepositoryError(err, "favorites not found")
	}

	return favorites, nil
}

func (r *gormFavoriteRepository) Exists(ctx context.Context, userID, carID uint) (bool, error) {
	var favorite domains.Favorite

	err := r.db.WithContext(ctx).
		Model(&domains.Favorite{}).
		Select("user_id", "car_id").
		Where("user_id = ? AND car_id = ?", userID, carID).
		Take(&favorite).Error
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	default:
		return false, mapRepositoryError(err, "favorite not found")
	}
}

package repository

import (
	"context"
	"errors"

	"car-service/internal/domains"
	apperrors "car-service/internal/errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type CarFilter struct {
	Search       string
	Brand        string
	Model        string
	YearFrom     *int
	YearTo       *int
	FuelType     string
	Transmission string
	BodyType     string
	SeatsMin     *int
	PriceMin     *int64
	PriceMax     *int64
	Purpose      string
	SortBy       string
	SortOrder    string
	Limit        int
	Offset       int
}

type CreateCarImageParams struct {
	CarID       uint
	BucketName  string
	ObjectKey   string
	FileName    string
	ContentType string
	FileSize    int64
	IsMain      bool
	SortOrder   int
}

type CarRepository interface {
	List(ctx context.Context, filter CarFilter) ([]domains.Car, error)
	Count(ctx context.Context, filter CarFilter) (int64, error)
	GetByID(ctx context.Context, id uint) (domains.Car, error)
	ExistsByID(ctx context.Context, id uint) (bool, error)
	CreateImage(ctx context.Context, params CreateCarImageParams) (domains.CarImage, error)
}

type gormCarRepository struct {
	db *gorm.DB
}

type carCatalogRow struct {
	ID              uint    `gorm:"column:id"`
	Brand           string  `gorm:"column:brand"`
	Model           string  `gorm:"column:model"`
	Year            int     `gorm:"column:year"`
	FuelType        string  `gorm:"column:fuel_type"`
	Transmission    string  `gorm:"column:transmission"`
	BodyType        string  `gorm:"column:body_type"`
	Color           string  `gorm:"column:color"`
	SeatsCount      int     `gorm:"column:seats_count"`
	PricePerDay     int64   `gorm:"column:price_per_day"`
	Purpose         string  `gorm:"column:purpose"`
	Description     string  `gorm:"column:description"`
	MainImageID     *uint   `gorm:"column:main_image_id"`
	MainImageBucket *string `gorm:"column:main_image_bucket_name"`
	MainImageKey    *string `gorm:"column:main_image_object_key"`
}

func NewCarRepository(db *gorm.DB) CarRepository {
	return &gormCarRepository{db: db}
}

func (r *gormCarRepository) List(ctx context.Context, filter CarFilter) ([]domains.Car, error) {
	var rows []carCatalogRow

	err := r.db.WithContext(ctx).
		Model(&domains.Car{}).
		Scopes(carCatalogFilterScopes(filter)...).
		Joins(
			"LEFT JOIN car_images AS main_image ON main_image.car_id = cars.id AND main_image.is_main = ?",
			true,
		).
		Select(carCatalogSelect()).
		Scopes(scopeCatalogSort(filter), scopeCatalogPagination(filter)).
		Find(&rows).Error
	if err != nil {
		return nil, mapRepositoryError(err, "cars not found")
	}

	cars := make([]domains.Car, 0, len(rows))
	for _, row := range rows {
		cars = append(cars, row.toDomain())
	}

	return cars, nil
}

func (r *gormCarRepository) Count(ctx context.Context, filter CarFilter) (int64, error) {
	var total int64

	err := r.db.WithContext(ctx).
		Model(&domains.Car{}).
		Scopes(carCatalogFilterScopes(filter)...).
		Count(&total).Error
	if err != nil {
		return 0, mapRepositoryError(err, "cars not found")
	}

	return total, nil
}

func (r *gormCarRepository) GetByID(ctx context.Context, id uint) (domains.Car, error) {
	var car domains.Car

	err := r.db.WithContext(ctx).
		Preload("CarImages", func(db *gorm.DB) *gorm.DB {
			return db.Order("is_main DESC").Order("sort_order ASC").Order("id ASC")
		}).
		Where("cars.id = ?", id).
		First(&car).Error
	if err != nil {
		return domains.Car{}, mapRepositoryError(err, "car not found")
	}

	return car, nil
}

func (r *gormCarRepository) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var car domains.Car

	err := r.db.WithContext(ctx).
		Model(&domains.Car{}).
		Select("id").
		Where("id = ?", id).
		Take(&car).Error
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	default:
		return false, mapRepositoryError(err, "car not found")
	}
}

func (r *gormCarRepository) CreateImage(
	ctx context.Context,
	params CreateCarImageParams,
) (domains.CarImage, error) {
	image := domains.CarImage{
		CarID:       params.CarID,
		BucketName:  params.BucketName,
		ObjectKey:   params.ObjectKey,
		FileName:    params.FileName,
		ContentType: params.ContentType,
		FileSize:    params.FileSize,
		IsMain:      params.IsMain,
		SortOrder:   params.SortOrder,
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if params.IsMain {
			if err := tx.Model(&domains.CarImage{}).
				Where("car_id = ? AND is_main = ?", params.CarID, true).
				Update("is_main", false).Error; err != nil {
				return err
			}
		}

		return tx.Create(&image).Error
	})
	if err != nil {
		return domains.CarImage{}, mapRepositoryError(err, "car image not saved")
	}

	return image, nil
}

func carCatalogSelect() []string {
	return []string{
		"cars.id",
		"cars.brand",
		"cars.model",
		"cars.year",
		"cars.fuel_type",
		"cars.transmission",
		"cars.body_type",
		"cars.color",
		"cars.seats_count",
		"cars.price_per_day",
		"cars.purpose",
		"cars.description",
		"main_image.id AS main_image_id",
		"main_image.bucket_name AS main_image_bucket_name",
		"main_image.object_key AS main_image_object_key",
	}
}

func (r carCatalogRow) toDomain() domains.Car {
	car := domains.Car{
		ID:           r.ID,
		Brand:        r.Brand,
		Model:        r.Model,
		Year:         r.Year,
		FuelType:     r.FuelType,
		Transmission: r.Transmission,
		BodyType:     r.BodyType,
		Color:        r.Color,
		SeatsCount:   r.SeatsCount,
		PricePerDay:  r.PricePerDay,
		Purpose:      r.Purpose,
		Description:  r.Description,
	}

	if r.MainImageID != nil && r.MainImageBucket != nil && r.MainImageKey != nil {
		car.CarImages = []domains.CarImage{
			{
				ID:         *r.MainImageID,
				CarID:      r.ID,
				BucketName: *r.MainImageBucket,
				ObjectKey:  *r.MainImageKey,
				IsMain:     true,
			},
		}
	}

	return car
}

func mapRepositoryError(err error, notFoundMessage string) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return apperrors.New(apperrors.ErrNotFound, notFoundMessage)
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return apperrors.New(apperrors.ErrConflict, "resource already exists")
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return apperrors.New(apperrors.ErrNotFound, "related resource not found")
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return apperrors.New(apperrors.ErrConflict, "resource already exists")
		case "23503":
			return apperrors.New(apperrors.ErrNotFound, "related resource not found")
		}
	}

	return err
}

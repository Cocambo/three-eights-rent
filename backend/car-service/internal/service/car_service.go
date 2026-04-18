package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"car-service/internal/domains"
	apperrors "car-service/internal/errors"
	"car-service/internal/repository"
	"car-service/internal/storage"
)

const (
	defaultCatalogLimit     = 20
	maxCatalogLimit         = 100
	defaultCatalogSortBy    = "id"
	defaultCatalogSortOrder = "desc"
)

var allowedCatalogSortFields = map[string]struct{}{
	"id":            {},
	"year":          {},
	"price_per_day": {},
	"created_at":    {},
}

type CatalogFilter struct {
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
	Limit        *int
	Offset       *int
}

type CatalogCar struct {
	ID           uint
	Brand        string
	Model        string
	Year         int
	FuelType     string
	Transmission string
	BodyType     string
	SeatsCount   int
	PricePerDay  int64
	Purpose      string
	MainImageURL *string
}

type CatalogResult struct {
	Items  []CatalogCar
	Total  int64
	Limit  int
	Offset int
}

type CarImageResult struct {
	ID        uint
	URL       string
	IsMain    bool
	SortOrder int
}

type CarDetailsResult struct {
	ID           uint
	Brand        string
	Model        string
	Year         int
	FuelType     string
	Transmission string
	BodyType     string
	Color        string
	SeatsCount   int
	PricePerDay  int64
	Purpose      string
	Description  string
	Images       []CarImageResult
}

type CarService interface {
	GetCatalog(ctx context.Context, filter CatalogFilter) (CatalogResult, error)
	GetCarDetails(ctx context.Context, carID uint) (CarDetailsResult, error)
}

type carService struct {
	carRepository repository.CarRepository
	imageStorage  storage.ImageStorageService
}

type normalizedCatalogFilter struct {
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

func NewCarService(
	carRepository repository.CarRepository,
	imageStorage storage.ImageStorageService,
) CarService {
	return &carService{
		carRepository: carRepository,
		imageStorage:  imageStorage,
	}
}

func (s *carService) GetCatalog(
	ctx context.Context,
	filter CatalogFilter,
) (CatalogResult, error) {
	normalized, err := normalizeCatalogFilter(filter)
	if err != nil {
		return CatalogResult{}, err
	}

	repositoryFilter := repository.CarFilter{
		Search:       normalized.Search,
		Brand:        normalized.Brand,
		Model:        normalized.Model,
		YearFrom:     normalized.YearFrom,
		YearTo:       normalized.YearTo,
		FuelType:     normalized.FuelType,
		Transmission: normalized.Transmission,
		BodyType:     normalized.BodyType,
		SeatsMin:     normalized.SeatsMin,
		PriceMin:     normalized.PriceMin,
		PriceMax:     normalized.PriceMax,
		Purpose:      normalized.Purpose,
		SortBy:       normalized.SortBy,
		SortOrder:    normalized.SortOrder,
		Limit:        normalized.Limit,
		Offset:       normalized.Offset,
	}

	total, err := s.carRepository.Count(ctx, repositoryFilter)
	if err != nil {
		return CatalogResult{}, err
	}

	cars, err := s.carRepository.List(ctx, repositoryFilter)
	if err != nil {
		return CatalogResult{}, err
	}

	items := make([]CatalogCar, 0, len(cars))
	for i := range cars {
		cars[i].CarImages = catalogImages(cars[i].CarImages)

		item, err := toCatalogCar(ctx, s.imageStorage, cars[i])
		if err != nil {
			return CatalogResult{}, err
		}

		items = append(items, item)
	}

	return CatalogResult{
		Items:  items,
		Total:  total,
		Limit:  normalized.Limit,
		Offset: normalized.Offset,
	}, nil
}

func (s *carService) GetCarDetails(ctx context.Context, carID uint) (CarDetailsResult, error) {
	if carID == 0 {
		return CarDetailsResult{}, validationError("car_id must be greater than zero")
	}

	car, err := s.carRepository.GetByID(ctx, carID)
	if err != nil {
		return CarDetailsResult{}, err
	}

	sortCarImages(car.CarImages)

	return toCarDetails(ctx, s.imageStorage, car)
}

func normalizeCatalogFilter(filter CatalogFilter) (normalizedCatalogFilter, error) {
	filter = trimCatalogFilter(filter)

	if filter.YearFrom != nil && filter.YearTo != nil && *filter.YearFrom > *filter.YearTo {
		return normalizedCatalogFilter{}, validationError("year_to must be greater than or equal to year_from")
	}

	if filter.PriceMin != nil && filter.PriceMax != nil && *filter.PriceMin > *filter.PriceMax {
		return normalizedCatalogFilter{}, validationError("price_max must be greater than or equal to price_min")
	}

	limit, err := resolveCatalogLimit(filter.Limit)
	if err != nil {
		return normalizedCatalogFilter{}, err
	}

	offset, err := resolveCatalogOffset(filter.Offset)
	if err != nil {
		return normalizedCatalogFilter{}, err
	}

	sortBy, sortOrder, err := resolveCatalogSorting(filter.SortBy, filter.SortOrder)
	if err != nil {
		return normalizedCatalogFilter{}, err
	}

	return normalizedCatalogFilter{
		Search:       filter.Search,
		Brand:        filter.Brand,
		Model:        filter.Model,
		YearFrom:     filter.YearFrom,
		YearTo:       filter.YearTo,
		FuelType:     filter.FuelType,
		Transmission: filter.Transmission,
		BodyType:     filter.BodyType,
		SeatsMin:     filter.SeatsMin,
		PriceMin:     filter.PriceMin,
		PriceMax:     filter.PriceMax,
		Purpose:      filter.Purpose,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
		Limit:        limit,
		Offset:       offset,
	}, nil
}

func trimCatalogFilter(filter CatalogFilter) CatalogFilter {
	filter.Search = strings.TrimSpace(filter.Search)
	filter.Brand = strings.TrimSpace(filter.Brand)
	filter.Model = strings.TrimSpace(filter.Model)
	filter.FuelType = strings.TrimSpace(filter.FuelType)
	filter.Transmission = strings.TrimSpace(filter.Transmission)
	filter.BodyType = strings.TrimSpace(filter.BodyType)
	filter.Purpose = strings.TrimSpace(filter.Purpose)
	filter.SortBy = strings.TrimSpace(filter.SortBy)
	filter.SortOrder = strings.ToLower(strings.TrimSpace(filter.SortOrder))

	return filter
}

func resolveCatalogLimit(limit *int) (int, error) {
	if limit == nil {
		return defaultCatalogLimit, nil
	}

	if *limit <= 0 {
		return 0, validationError("limit must be greater than zero")
	}

	if *limit > maxCatalogLimit {
		return maxCatalogLimit, nil
	}

	return *limit, nil
}

func resolveCatalogOffset(offset *int) (int, error) {
	if offset == nil {
		return 0, nil
	}

	if *offset < 0 {
		return 0, validationError("offset must be greater than or equal to zero")
	}

	return *offset, nil
}

func resolveCatalogSorting(sortBy, sortOrder string) (string, string, error) {
	if sortBy == "" {
		sortBy = defaultCatalogSortBy
	}

	if _, ok := allowedCatalogSortFields[sortBy]; !ok {
		return "", "", validationError("sort_by has invalid value")
	}

	if sortOrder == "" {
		sortOrder = defaultCatalogSortOrder
	}

	switch sortOrder {
	case "asc", "desc":
		return sortBy, sortOrder, nil
	default:
		return "", "", validationError("sort_order has invalid value")
	}
}

func catalogImages(images []domains.CarImage) []domains.CarImage {
	for _, image := range images {
		if image.IsMain {
			return []domains.CarImage{image}
		}
	}

	return nil
}

func sortCarImages(images []domains.CarImage) {
	sort.SliceStable(images, func(i, j int) bool {
		left := images[i]
		right := images[j]

		if left.IsMain != right.IsMain {
			return left.IsMain
		}

		if left.SortOrder != right.SortOrder {
			return left.SortOrder < right.SortOrder
		}

		return left.ID < right.ID
	})
}

func toCatalogCar(
	ctx context.Context,
	imageStorage storage.ImageStorageService,
	car domains.Car,
) (CatalogCar, error) {
	mainImageURL, err := mainImageURL(ctx, imageStorage, car.CarImages)
	if err != nil {
		return CatalogCar{}, err
	}

	return CatalogCar{
		ID:           car.ID,
		Brand:        car.Brand,
		Model:        car.Model,
		Year:         car.Year,
		FuelType:     car.FuelType,
		Transmission: car.Transmission,
		BodyType:     car.BodyType,
		SeatsCount:   car.SeatsCount,
		PricePerDay:  car.PricePerDay,
		Purpose:      car.Purpose,
		MainImageURL: mainImageURL,
	}, nil
}

func toCarDetails(
	ctx context.Context,
	imageStorage storage.ImageStorageService,
	car domains.Car,
) (CarDetailsResult, error) {
	images := make([]CarImageResult, 0, len(car.CarImages))
	for _, image := range car.CarImages {
		presignedURL, err := imageStorage.GeneratePresignedGetURL(ctx, image.BucketName, image.ObjectKey)
		if err != nil {
			return CarDetailsResult{}, fmt.Errorf(
				"generate image url for car %d image %d: %w",
				car.ID,
				image.ID,
				err,
			)
		}

		images = append(images, CarImageResult{
			ID:        image.ID,
			URL:       presignedURL,
			IsMain:    image.IsMain,
			SortOrder: image.SortOrder,
		})
	}

	return CarDetailsResult{
		ID:           car.ID,
		Brand:        car.Brand,
		Model:        car.Model,
		Year:         car.Year,
		FuelType:     car.FuelType,
		Transmission: car.Transmission,
		BodyType:     car.BodyType,
		Color:        car.Color,
		SeatsCount:   car.SeatsCount,
		PricePerDay:  car.PricePerDay,
		Purpose:      car.Purpose,
		Description:  car.Description,
		Images:       images,
	}, nil
}

func mainImageURL(
	ctx context.Context,
	imageStorage storage.ImageStorageService,
	images []domains.CarImage,
) (*string, error) {
	if len(images) == 0 {
		return nil, nil
	}

	presignedURL, err := imageStorage.GeneratePresignedGetURL(
		ctx,
		images[0].BucketName,
		images[0].ObjectKey,
	)
	if err != nil {
		return nil, fmt.Errorf("generate main image url for image %d: %w", images[0].ID, err)
	}

	return &presignedURL, nil
}

func validationError(message string) error {
	return apperrors.New(apperrors.ErrValidation, message)
}

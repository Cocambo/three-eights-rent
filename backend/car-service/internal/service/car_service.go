package service

import (
	"context"
	"sort"
	"strings"

	"car-service/internal/domains"
	apperrors "car-service/internal/errors"
	"car-service/internal/repository"
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

type Catalog struct {
	Items  []domains.Car
	Total  int64
	Limit  int
	Offset int
}

type CarService interface {
	GetCatalog(ctx context.Context, filter CatalogFilter) (Catalog, error)
	GetCarDetails(ctx context.Context, carID uint) (domains.Car, error)
}

type carService struct {
	carRepository repository.CarRepository
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

func NewCarService(carRepository repository.CarRepository) CarService {
	return &carService{carRepository: carRepository}
}

func (s *carService) GetCatalog(ctx context.Context, filter CatalogFilter) (Catalog, error) {
	normalized, err := normalizeCatalogFilter(filter)
	if err != nil {
		return Catalog{}, err
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
		return Catalog{}, err
	}

	cars, err := s.carRepository.List(ctx, repositoryFilter)
	if err != nil {
		return Catalog{}, err
	}

	for i := range cars {
		cars[i].CarImages = catalogImages(cars[i].CarImages)
	}

	return Catalog{
		Items:  cars,
		Total:  total,
		Limit:  normalized.Limit,
		Offset: normalized.Offset,
	}, nil
}

func (s *carService) GetCarDetails(ctx context.Context, carID uint) (domains.Car, error) {
	if carID == 0 {
		return domains.Car{}, validationError("car_id must be greater than zero")
	}

	car, err := s.carRepository.GetByID(ctx, carID)
	if err != nil {
		return domains.Car{}, err
	}

	sortCarImages(car.CarImages)

	return car, nil
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

func validationError(message string) error {
	return apperrors.New(apperrors.ErrValidation, message)
}

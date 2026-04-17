package handler

import (
	"net/http"

	"car-service/internal/domains"
	"car-service/internal/dto"
	apperrors "car-service/internal/errors"
	"car-service/internal/service"

	"github.com/gin-gonic/gin"
)

type CarHandler struct {
	carService service.CarService
}

func NewCarHandler(carService service.CarService) *CarHandler {
	return &CarHandler{carService: carService}
}

func (h *CarHandler) Health(c *gin.Context) {
	writeSuccess(c, http.StatusOK, dto.HealthResponse{
		Status: "ok",
	})
}

func (h *CarHandler) List(c *gin.Context) {
	var req dto.ListCarsQuery
	if !bindQuery(c, &req) {
		return
	}

	req.Normalize()
	if err := req.Validate(); err != nil {
		writeError(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	cars, err := h.carService.List(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}

	total := int64(len(cars))
	catalogItems := make([]dto.CarCatalogItemResponse, 0, len(cars))
	for _, car := range cars {
		catalogItems = append(catalogItems, toCarCatalogItemResponse(car))
	}

	items := paginateCatalogItems(catalogItems, req.OffsetValue(), req.LimitValue())

	writeSuccess(c, http.StatusOK, dto.CarsCatalogResponse{
		Items: items,
		Pagination: dto.PaginationMeta{
			Total:  total,
			Limit:  req.LimitValue(),
			Offset: req.OffsetValue(),
		},
	})
}

func (h *CarHandler) GetByID(c *gin.Context) {
	var req dto.GetCarByIDURI
	if !bindURI(c, &req) {
		return
	}

	car, err := h.carService.GetByID(c.Request.Context(), req.ID)
	if err != nil {
		writeError(c, err)
		return
	}

	writeSuccess(c, http.StatusOK, toCarDetailsResponse(car))
}

func toCarCatalogItemResponse(car domains.Car) dto.CarCatalogItemResponse {
	return dto.CarCatalogItemResponse{
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
		MainImageURL: mainImageURL(car.CarImages),
	}
}

func toCarDetailsResponse(car domains.Car) dto.CarDetailsResponse {
	images := make([]dto.CarImageResponse, 0, len(car.CarImages))
	for _, image := range car.CarImages {
		images = append(images, dto.CarImageResponse{
			ID:        image.ID,
			URL:       image.ObjectKey,
			IsMain:    image.IsMain,
			SortOrder: image.SortOrder,
		})
	}

	return dto.CarDetailsResponse{
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
	}
}

func mainImageURL(images []domains.CarImage) string {
	for _, image := range images {
		if image.IsMain {
			return image.ObjectKey
		}
	}

	if len(images) == 0 {
		return ""
	}

	return images[0].ObjectKey
}

func paginateCatalogItems(items []dto.CarCatalogItemResponse, offset, limit int) []dto.CarCatalogItemResponse {
	if offset >= len(items) {
		return []dto.CarCatalogItemResponse{}
	}

	end := offset + limit
	if end > len(items) {
		end = len(items)
	}

	return items[offset:end]
}

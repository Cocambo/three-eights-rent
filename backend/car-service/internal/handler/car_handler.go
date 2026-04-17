package handler

import (
	"net/http"

	"car-service/internal/domains"
	"car-service/internal/dto"
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

	catalog, err := h.carService.GetCatalog(c.Request.Context(), service.CatalogFilter{
		Search:       req.Q,
		Brand:        req.Brand,
		Model:        req.Model,
		YearFrom:     req.YearFrom,
		YearTo:       req.YearTo,
		FuelType:     req.FuelType,
		Transmission: req.Transmission,
		BodyType:     req.BodyType,
		SeatsMin:     req.SeatsMin,
		PriceMin:     req.PriceMin,
		PriceMax:     req.PriceMax,
		Purpose:      req.Purpose,
		SortBy:       req.SortBy,
		SortOrder:    req.SortOrder,
		Limit:        req.Limit,
		Offset:       req.Offset,
	})
	if err != nil {
		writeError(c, err)
		return
	}

	catalogItems := make([]dto.CarCatalogItemResponse, 0, len(catalog.Items))
	for _, car := range catalog.Items {
		catalogItems = append(catalogItems, toCarCatalogItemResponse(car))
	}

	writeSuccess(c, http.StatusOK, dto.CarsCatalogResponse{
		Items: catalogItems,
		Pagination: dto.PaginationMeta{
			Total:  catalog.Total,
			Limit:  catalog.Limit,
			Offset: catalog.Offset,
		},
	})
}

func (h *CarHandler) GetByID(c *gin.Context) {
	var req dto.GetCarByIDURI
	if !bindURI(c, &req) {
		return
	}

	car, err := h.carService.GetCarDetails(c.Request.Context(), req.ID)
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

func mainImageURL(images []domains.CarImage) *string {
	if len(images) == 0 {
		return nil
	}

	url := images[0].ObjectKey
	return &url
}

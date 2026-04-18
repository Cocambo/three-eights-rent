package handler

import (
	"net/http"

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

	writeSuccess(c, http.StatusOK, catalog)
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

	writeSuccess(c, http.StatusOK, car)
}

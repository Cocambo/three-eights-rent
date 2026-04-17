package handler

import (
	"net/http"
	"strconv"

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
	cars, err := h.carService.List(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}

	items := make([]dto.CarResponse, 0, len(cars))
	for _, car := range cars {
		items = append(items, toCarResponse(car))
	}

	writeSuccess(c, http.StatusOK, items)
}

func (h *CarHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeError(c, apperrors.New(apperrors.ErrValidation, "invalid car id"))
		return
	}

	car, err := h.carService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		writeError(c, err)
		return
	}

	writeSuccess(c, http.StatusOK, toCarResponse(car))
}

func toCarResponse(car domains.Car) dto.CarResponse {
	return dto.CarResponse{
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
	}
}

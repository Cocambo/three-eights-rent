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
	c.JSON(http.StatusOK, dto.HealthResponse{
		Status: "ok",
	})
}

func (h *CarHandler) List(c *gin.Context) {
	cars, err := h.carService.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch cars",
		})
		return
	}

	response := make([]dto.CarResponse, 0, len(cars))
	for _, car := range cars {
		response = append(response, dto.CarResponse{
			ID:             car.ID,
			Brand:          car.Brand,
			Model:          car.Model,
			Year:           car.Year,
			RegistrationNo: car.RegistrationNo,
			ImageURL:       car.ImageURL,
		})
	}

	c.JSON(http.StatusOK, response)
}

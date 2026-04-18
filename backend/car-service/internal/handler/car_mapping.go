package handler

import (
	"car-service/internal/domains"
	"car-service/internal/dto"
)

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

func mainImageURL(images []domains.CarImage) *string {
	if len(images) == 0 {
		return nil
	}

	url := images[0].ObjectKey
	return &url
}

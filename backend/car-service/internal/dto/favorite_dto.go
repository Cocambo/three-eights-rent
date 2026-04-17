package dto

import "time"

type ListFavoritesQuery struct{}

type FavoriteCarURI struct {
	CarID uint `uri:"carId" binding:"required,gt=0"`
}

type FavoriteResponse struct {
	CarID   uint                   `json:"car_id"`
	AddedAt time.Time              `json:"added_at"`
	Car     CarCatalogItemResponse `json:"car"`
}

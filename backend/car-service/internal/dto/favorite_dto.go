package dto

import "time"

type ListFavoritesQuery struct{}

type FavoriteCarURI struct {
	CarID uint `uri:"carId" binding:"required,gt=0"`
}

type ListFavoritesResponse struct {
	Items []FavoriteResponse `json:"items"`
}

type FavoriteResponse struct {
	CarID   uint                   `json:"car_id"`
	AddedAt time.Time              `json:"added_at"`
	Car     CarCatalogItemResponse `json:"car"`
}

type FavoriteMutationResponse struct {
	CarID   uint   `json:"car_id"`
	Message string `json:"message"`
}

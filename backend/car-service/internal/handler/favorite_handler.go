package handler

import (
	"net/http"

	"car-service/internal/dto"
	apperrors "car-service/internal/errors"
	"car-service/internal/middleware"
	"car-service/internal/service"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	favoriteService service.FavoriteService
}

func NewFavoriteHandler(favoriteService service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{favoriteService: favoriteService}
}

func (h *FavoriteHandler) List(c *gin.Context) {
	userID, ok := middleware.UserIDFromGin(c)
	if !ok {
		writeError(c, serviceUnauthorizedError())
		return
	}

	favorites, err := h.favoriteService.GetUserFavorites(c.Request.Context(), userID)
	if err != nil {
		writeError(c, err)
		return
	}

	response := make([]dto.FavoriteResponse, 0, len(favorites))
	for _, favorite := range favorites {
		response = append(response, dto.FavoriteResponse{
			CarID:   favorite.CarID,
			AddedAt: favorite.CreatedAt,
			Car:     toCarCatalogItemResponse(favorite.Car),
		})
	}

	writeSuccess(c, http.StatusOK, response)
}

func (h *FavoriteHandler) Add(c *gin.Context) {
	userID, ok := middleware.UserIDFromGin(c)
	if !ok {
		writeError(c, serviceUnauthorizedError())
		return
	}

	var req dto.FavoriteCarURI
	if !bindURI(c, &req) {
		return
	}

	if err := h.favoriteService.AddToFavorites(c.Request.Context(), userID, req.CarID); err != nil {
		writeError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *FavoriteHandler) Remove(c *gin.Context) {
	userID, ok := middleware.UserIDFromGin(c)
	if !ok {
		writeError(c, serviceUnauthorizedError())
		return
	}

	var req dto.FavoriteCarURI
	if !bindURI(c, &req) {
		return
	}

	if err := h.favoriteService.RemoveFromFavorites(c.Request.Context(), userID, req.CarID); err != nil {
		writeError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func serviceUnauthorizedError() error {
	return apperrors.New(apperrors.ErrUnauthorized, "unauthorized")
}

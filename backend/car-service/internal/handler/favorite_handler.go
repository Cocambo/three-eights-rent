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
	var req dto.ListFavoritesQuery
	if !bindQuery(c, &req) {
		return
	}

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

	writeSuccess(c, http.StatusOK, toListFavoritesResponse(favorites))
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

	writeSuccess(c, http.StatusOK, dto.FavoriteMutationResponse{
		CarID:   req.CarID,
		Message: "car added to favorites",
	})
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

	writeSuccess(c, http.StatusOK, dto.FavoriteMutationResponse{
		CarID:   req.CarID,
		Message: "car removed from favorites",
	})
}

func serviceUnauthorizedError() error {
	return apperrors.New(apperrors.ErrUnauthorized, "unauthorized")
}

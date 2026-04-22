package handler

import (
	"net/http"

	"car-service/internal/middleware"
	"car-service/internal/service"

	"github.com/gin-gonic/gin"
)

type RecommendationHandler struct {
	recommendationService service.RecommendationService
}

func NewRecommendationHandler(recommendationService service.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{recommendationService: recommendationService}
}

func (h *RecommendationHandler) GetMe(c *gin.Context) {
	userID, ok := middleware.UserIDFromGin(c)
	if !ok {
		writeError(c, serviceUnauthorizedError())
		return
	}

	recommendations, err := h.recommendationService.GetMyRecommendations(c.Request.Context(), userID)
	if err != nil {
		writeError(c, err)
		return
	}

	writeSuccess(c, http.StatusOK, toRecommendationsResponse(recommendations))
}

func (h *RecommendationHandler) Rebuild(c *gin.Context) {
	result, err := h.recommendationService.RebuildRecommendations(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}

	writeSuccess(c, http.StatusOK, toRebuildRecommendationsResponse(result))
}

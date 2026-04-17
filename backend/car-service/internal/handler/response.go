package handler

import (
	"errors"
	"net/http"

	apperrors "car-service/internal/errors"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func MapError(err error) (int, string) {
	switch {
	case errors.Is(err, apperrors.ErrValidation):
		return http.StatusBadRequest, "validation"
	case errors.Is(err, apperrors.ErrUnauthorized):
		return http.StatusUnauthorized, "unauthorized"
	case errors.Is(err, apperrors.ErrNotFound):
		return http.StatusNotFound, "not_found"
	case errors.Is(err, apperrors.ErrConflict):
		return http.StatusConflict, "conflict"
	default:
		return http.StatusInternalServerError, "internal_error"
	}
}

func writeError(c *gin.Context, err error) {
	status, code := MapError(err)
	message := err.Error()
	if status == http.StatusInternalServerError {
		message = "internal server error"
	}

	c.AbortWithStatusJSON(status, ErrorResponse{
		Error: ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}

func writeSuccess(c *gin.Context, status int, data any) {
	c.JSON(status, data)
}

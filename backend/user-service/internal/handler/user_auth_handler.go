package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"user-service/internal/dto"
	apperrors "user-service/internal/errors"
	"user-service/internal/middleware"
	"user-service/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if !bindJSON(c, &req) {
		return
	}

	resp, err := h.userService.Register(c.Request.Context(), req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if !bindJSON(c, &req) {
		return
	}

	resp, err := h.userService.Login(c.Request.Context(), req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if !bindJSON(c, &req) {
		return
	}

	resp, err := h.userService.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if !bindJSON(c, &req) {
		return
	}

	if err := h.userService.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		handleServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, ok := middleware.UserIDFromContext(c.Request.Context())
	if !ok {
		writeError(c, http.StatusUnauthorized, "user_id is missing in context")
		return
	}

	resp, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.UserIDFromContext(c.Request.Context())
	if !ok {
		writeError(c, http.StatusUnauthorized, "user_id is missing in context")
		return
	}

	var req dto.UpdateProfileRequest
	if !bindJSON(c, &req) {
		return
	}

	resp, err := h.userService.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) CreateDriverLicense(c *gin.Context) {
	userID, ok := middleware.UserIDFromContext(c.Request.Context())
	if !ok {
		writeError(c, http.StatusUnauthorized, "user_id is missing in context")
		return
	}

	var req dto.CreateDriverLicenseRequest
	if !bindJSON(c, &req) {
		return
	}

	resp, err := h.userService.CreateDriverLicense(c.Request.Context(), userID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperrors.ErrValidation):
		writeError(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, apperrors.ErrDuplicateEmail), errors.Is(err, apperrors.ErrConflict):
		writeError(c, http.StatusConflict, err.Error())
	case errors.Is(err, apperrors.ErrInvalidCredentials):
		writeError(c, http.StatusUnauthorized, err.Error())
	case errors.Is(err, apperrors.ErrInvalidToken), errors.Is(err, apperrors.ErrUnauthorized):
		writeError(c, http.StatusUnauthorized, err.Error())
	case errors.Is(err, apperrors.ErrNotFound):
		writeError(c, http.StatusNotFound, err.Error())
	default:
		writeError(c, http.StatusInternalServerError, "internal server error")
	}
}

func writeError(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{
		"error": message,
	})
}

func bindJSON(c *gin.Context, dst any) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		writeError(c, http.StatusBadRequest, formatBindingError(err))
		return false
	}
	return true
}

func formatBindingError(err error) string {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) && len(validationErrs) > 0 {
		validationErr := validationErrs[0]
		field := validationErr.Field()
		switch validationErr.Tag() {
		case "required":
			return field + " is required"
		case "email":
			return field + " must be a valid email"
		case "min":
			return field + " is too short"
		case "max":
			return field + " is too long"
		case "datetime":
			return field + " must match YYYY-MM-DD"
		default:
			return field + " is invalid"
		}
	}
	return err.Error()
}

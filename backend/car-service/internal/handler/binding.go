package handler

import (
	"errors"
	"strings"

	apperrors "car-service/internal/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func bindQuery(c *gin.Context, dst any) bool {
	if err := c.ShouldBindQuery(dst); err != nil {
		writeError(c, apperrors.New(apperrors.ErrValidation, formatBindingError(err)))
		return false
	}

	return true
}

func bindURI(c *gin.Context, dst any) bool {
	if err := c.ShouldBindUri(dst); err != nil {
		writeError(c, apperrors.New(apperrors.ErrValidation, formatBindingError(err)))
		return false
	}

	return true
}

func bindForm(c *gin.Context, dst any) bool {
	if err := c.ShouldBind(dst); err != nil {
		writeError(c, apperrors.New(apperrors.ErrValidation, formatBindingError(err)))
		return false
	}

	return true
}

func formatBindingError(err error) string {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) && len(validationErrs) > 0 {
		validationErr := validationErrs[0]
		field := toSnakeCase(validationErr.Field())

		switch validationErr.Tag() {
		case "required":
			return field + " is required"
		case "oneof":
			return field + " has invalid value"
		case "gt", "gte":
			return field + " is too small"
		case "lt", "lte":
			return field + " is too large"
		case "max":
			return field + " is too long"
		default:
			return field + " is invalid"
		}
	}

	if strings.Contains(err.Error(), "http: request body too large") {
		return "request body is too large"
	}

	return err.Error()
}

func toSnakeCase(value string) string {
	var builder strings.Builder

	for i, r := range value {
		if i > 0 && r >= 'A' && r <= 'Z' {
			prev := rune(value[i-1])
			nextIsLower := i+1 < len(value) && value[i+1] >= 'a' && value[i+1] <= 'z'
			if (prev >= 'a' && prev <= 'z') || nextIsLower {
				builder.WriteByte('_')
			}
		}

		if r >= 'A' && r <= 'Z' {
			builder.WriteByte(byte(r + ('a' - 'A')))
			continue
		}

		builder.WriteRune(r)
	}

	return builder.String()
}

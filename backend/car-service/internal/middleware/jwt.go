package middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"

	apperrors "car-service/internal/errors"
	httpresponse "car-service/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	userIDGinContextKey     = "user_id"
	userIDRequestContextKey = contextKey("user_id")
)

type JWTMiddleware struct {
	accessSecret []byte
}

type AccessTokenClaims struct {
	TokenType string `json:"token_type"`
	UserID    uint   `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTMiddleware(accessSecret string) *JWTMiddleware {
	return &JWTMiddleware{
		accessSecret: []byte(accessSecret),
	}
}

func (m *JWTMiddleware) JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := extractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			abortUnauthorized(c, err.Error())
			return
		}

		claims, err := m.parseAccessToken(tokenString)
		if err != nil {
			abortUnauthorized(c, err.Error())
			return
		}

		ctx := context.WithValue(c.Request.Context(), userIDRequestContextKey, claims.UserID)
		c.Request = c.Request.WithContext(ctx)
		c.Set(userIDGinContextKey, claims.UserID)
		c.Next()
	}
}

func (m *JWTMiddleware) parseAccessToken(tokenString string) (*AccessTokenClaims, error) {
	claims := &AccessTokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}

		return m.accessSecret, nil
	}, jwt.WithExpirationRequired(), jwt.WithIssuedAt())
	if err != nil {
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.TokenType != "access" {
		return nil, errors.New("invalid token")
	}

	if claims.UserID == 0 {
		return nil, errors.New("invalid token")
	}

	if claims.Subject != "" && claims.Subject != fmt.Sprintf("%d", claims.UserID) {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func UserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(userIDRequestContextKey).(uint)
	return userID, ok && userID > 0
}

func UserIDFromGin(c *gin.Context) (uint, bool) {
	userID, ok := c.Get(userIDGinContextKey)
	if !ok {
		return 0, false
	}

	typedUserID, ok := userID.(uint)
	return typedUserID, ok && typedUserID > 0
}

func extractBearerToken(header string) (string, error) {
	header = strings.TrimSpace(header)
	if header == "" {
		return "", errors.New("authorization header is required")
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("authorization header must use Bearer scheme")
	}

	tokenString := strings.TrimSpace(parts[1])
	if tokenString == "" {
		return "", errors.New("bearer token is required")
	}

	return tokenString, nil
}

func abortUnauthorized(c *gin.Context, message string) {
	httpresponse.WriteError(c, apperrors.New(apperrors.ErrUnauthorized, message))
}

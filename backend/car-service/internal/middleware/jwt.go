package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDContextKey contextKey = "user_id"

type JWTMiddleware struct {
	accessSecret string
}

type tokenClaims struct {
	TokenType string `json:"token_type"`
	UserID    uint   `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTMiddleware(accessSecret string) *JWTMiddleware {
	return &JWTMiddleware{accessSecret: accessSecret}
}

func (m *JWTMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := extractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		claims := &tokenClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(m.accessSecret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		if claims.TokenType != "access" || claims.UserID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		if claims.Subject != "" && claims.Subject != fmt.Sprintf("%d", claims.UserID) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		ctx := context.WithValue(c.Request.Context(), userIDContextKey, claims.UserID)
		c.Request = c.Request.WithContext(ctx)
		c.Set(string(userIDContextKey), claims.UserID)

		c.Next()
	}
}

func UserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(userIDContextKey).(uint)
	return userID, ok && userID > 0
}

func UserIDFromGin(c *gin.Context) (uint, bool) {
	userID, ok := c.Get(string(userIDContextKey))
	if !ok {
		return 0, false
	}

	typedUserID, ok := userID.(uint)
	return typedUserID, ok && typedUserID > 0
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("authorization header is required")
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("authorization header must use Bearer scheme")
	}

	return parts[1], nil
}

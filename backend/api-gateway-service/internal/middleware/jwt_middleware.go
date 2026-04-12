package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	HeaderUserID = "X-User-ID"
)

type TokenClaims struct {
	TokenType string `json:"token_type"`
	UserID    uint   `json:"user_id"`
	jwt.RegisteredClaims
}

func StripIdentityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		clearIdentityHeaders(c)
		c.Next()
	}
}

func JWTMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractBearerToken(c.GetHeader("Authorization"))
		if tokenString == "" {
			abortUnauthorized(c, "missing authentication token", "MISSING_TOKEN")
			return
		}

		claims := &TokenClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			abortUnauthorized(c, "invalid or expired token", "INVALID_TOKEN")
			return
		}

		if claims.TokenType != "access" || claims.UserID == 0 {
			abortUnauthorized(c, "invalid or expired token", "INVALID_TOKEN")
			return
		}
		if claims.Subject != "" && claims.Subject != fmt.Sprintf("%d", claims.UserID) {
			abortUnauthorized(c, "invalid or expired token", "INVALID_TOKEN")
			return
		}

		clearIdentityHeaders(c)
		c.Request.Header.Set(HeaderUserID, fmt.Sprintf("%d", claims.UserID))

		c.Next()
	}
}

func clearIdentityHeaders(c *gin.Context) {
	c.Request.Header.Del(HeaderUserID)
}

func extractBearerToken(header string) string {
	header = strings.TrimSpace(header)
	if header == "" {
		return ""
	}

	tokenType, token, ok := strings.Cut(header, " ")
	if !ok || !strings.EqualFold(tokenType, "Bearer") || strings.TrimSpace(token) == "" {
		return ""
	}

	return strings.TrimSpace(token)
}

func abortUnauthorized(c *gin.Context, message, code string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": message,
		"code":  code,
	})
}

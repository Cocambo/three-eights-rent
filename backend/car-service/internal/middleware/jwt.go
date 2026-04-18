package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const userIDContextKey = "user_id"

type JWTMiddleware struct {
	accessSecret []byte
}

type AccessTokenClaims struct {
	TokenType string `json:"token_type"`
	UserID    uint   `json:"user_id"`
	jwt.RegisteredClaims
}

type tokenErrorResponse struct {
	Error tokenErrorBody `json:"error"`
}

type tokenErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
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

		c.Set(userIDContextKey, claims.UserID)
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

func UserIDFromGin(c *gin.Context) (uint, bool) {
	userID, ok := c.Get(userIDContextKey)
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
	c.AbortWithStatusJSON(http.StatusUnauthorized, tokenErrorResponse{
		Error: tokenErrorBody{
			Code:    "unauthorized",
			Message: message,
		},
	})
}

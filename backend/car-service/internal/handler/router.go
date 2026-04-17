package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JWTMiddleware interface {
	Handler() gin.HandlerFunc
}

type Dependencies struct {
	CarHandler     *CarHandler
	JWTMiddleware  JWTMiddleware
	MinIOAvailable bool
}

func RegisterRoutes(router *gin.Engine, deps Dependencies) {
	router.GET("/health", deps.CarHandler.Health)

	api := router.Group("/api/v1")
	{
		api.GET("/info", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"service":          "car-service",
				"jwt_enabled":      deps.JWTMiddleware != nil,
				"minio_configured": deps.MinIOAvailable,
			})
		})
	}

	protected := router.Group("/api/v1")
	if deps.JWTMiddleware != nil {
		protected.Use(deps.JWTMiddleware.Handler())
	}
	{
		protected.GET("/cars", deps.CarHandler.List)
		protected.GET("/cars/:id", deps.CarHandler.GetByID)
	}
}

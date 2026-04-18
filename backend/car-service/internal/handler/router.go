package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JWTMiddleware interface {
	JWTAuthMiddleware() gin.HandlerFunc
}

type Dependencies struct {
	CarHandler      *CarHandler
	FavoriteHandler *FavoriteHandler
	JWTMiddleware   JWTMiddleware
	MinIOAvailable  bool
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

	public := api.Group("")
	{
		public.GET("/cars", deps.CarHandler.List)
		public.GET("/cars/:id", deps.CarHandler.GetByID)
	}

	favorites := api.Group("/favorites")
	if deps.JWTMiddleware != nil {
		favorites.Use(deps.JWTMiddleware.JWTAuthMiddleware())
	}
	{
		favorites.GET("", deps.FavoriteHandler.List)
		favorites.POST("/:carId", deps.FavoriteHandler.Add)
		favorites.DELETE("/:carId", deps.FavoriteHandler.Remove)
	}
}

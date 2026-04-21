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
	BookingHandler  *BookingHandler
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
		public.GET("/cars/:id/availability", deps.BookingHandler.GetAvailability)
		public.POST("/cars/:id/images", deps.CarHandler.UploadImage)
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

	bookings := api.Group("/bookings")
	if deps.JWTMiddleware != nil {
		bookings.Use(deps.JWTMiddleware.JWTAuthMiddleware())
	}
	{
		bookings.POST("", deps.BookingHandler.Create)
		bookings.DELETE("/:id", deps.BookingHandler.Cancel)
		bookings.GET("", deps.BookingHandler.List)
	}
}

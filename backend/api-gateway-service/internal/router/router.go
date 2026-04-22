package router

import (
	"net/url"

	"github.com/gin-gonic/gin"

	"api-gateway-service/internal/config"
	"api-gateway-service/internal/middleware"
	"api-gateway-service/internal/proxy"
)

func Register(engine *gin.Engine, cfg *config.Config) {
	userServiceURL, err := url.Parse(cfg.UserServiceURL)
	if err != nil {
		panic("validated user service url should always parse")
	}

	carServiceURL, err := url.Parse(cfg.CarServiceURL)
	if err != nil {
		panic("validated car service url should always parse")
	}

	userProxy := proxy.NewReverseProxy(userServiceURL, "/api/v1/users", "/api/v1")
	carProxy := proxy.NewReverseProxy(carServiceURL, "/api/v1/cars", "/api/v1/cars")
	favoritesProxy := proxy.NewReverseProxy(carServiceURL, "/api/v1/favorites", "/api/v1/favorites")
	bookingsProxy := proxy.NewReverseProxy(carServiceURL, "/api/v1/bookings", "/api/v1/bookings")
	recommendationsProxy := proxy.NewReverseProxy(carServiceURL, "/api/v1/recommendations", "/api/v1/recommendations")
	internalRecommendationsProxy := proxy.NewReverseProxy(carServiceURL, "/internal/recommendations", "/internal/recommendations")
	jwtMiddleware := middleware.JWTMiddleware(cfg.JWTAccessSecret)
	stripIdentityHeaders := middleware.StripIdentityHeaders()

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	publicUsers := engine.Group("/api/v1/users")
	publicUsers.Use(stripIdentityHeaders)
	{
		publicUsers.POST("/register", userProxy)
		publicUsers.POST("/login", userProxy)
		publicUsers.POST("/refresh", userProxy)
		publicUsers.POST("/logout", userProxy)
	}

	publicCars := engine.Group("/api/v1/cars")
	publicCars.Use(stripIdentityHeaders)
	{
		publicCars.GET("", carProxy)
		publicCars.GET("/:id", carProxy)
		publicCars.GET("/:id/availability", carProxy)
		publicCars.POST("/:id/images", carProxy)
	}

	protectedUsers := engine.Group("/api/v1/users")
	protectedUsers.Use(stripIdentityHeaders, jwtMiddleware)
	{
		protectedUsers.GET("/profile", userProxy)
		protectedUsers.PUT("/profile", userProxy)
		protectedUsers.GET("/driver-license", userProxy)
		protectedUsers.POST("/driver-license", userProxy)
	}

	protectedFavorites := engine.Group("/api/v1/favorites")
	protectedFavorites.Use(stripIdentityHeaders, jwtMiddleware)
	{
		protectedFavorites.GET("", favoritesProxy)
		protectedFavorites.POST("/:carId", favoritesProxy)
		protectedFavorites.DELETE("/:carId", favoritesProxy)
	}

	protectedBookings := engine.Group("/api/v1/bookings")
	protectedBookings.Use(stripIdentityHeaders, jwtMiddleware)
	{
		protectedBookings.POST("", bookingsProxy)
		protectedBookings.DELETE("/:id", bookingsProxy)
		protectedBookings.GET("", bookingsProxy)
	}

	protectedRecommendations := engine.Group("/api/v1/recommendations")
	protectedRecommendations.Use(stripIdentityHeaders, jwtMiddleware)
	{
		protectedRecommendations.GET("/me", recommendationsProxy)
	}

	protectedInternalRecommendations := engine.Group("/internal/recommendations")
	protectedInternalRecommendations.Use(stripIdentityHeaders, jwtMiddleware)
	{
		protectedInternalRecommendations.POST("/rebuild", internalRecommendationsProxy)
	}
}

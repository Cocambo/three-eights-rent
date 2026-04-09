package handler

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, userHandler *UserHandler, authMiddleware gin.HandlerFunc) {
	api := router.Group("/api/v1")
	{
		api.GET("/health", userHandler.HealthCheck)
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		api.POST("/refresh", userHandler.Refresh)
		api.POST("/logout", userHandler.Logout)

		protected := api.Group("")
		protected.Use(authMiddleware)
		{
			protected.GET("/profile", userHandler.GetProfile)
			protected.PUT("/profile", userHandler.UpdateProfile)
			protected.GET("/driver-license", userHandler.GetDriverLicense)
			protected.POST("/driver-license", userHandler.CreateDriverLicense)
		}
	}
}

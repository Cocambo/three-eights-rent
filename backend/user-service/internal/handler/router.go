package handler

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, userHandler *UserHandler) {
	api := router.Group("/api/v1")
	{
		api.GET("/health", userHandler.HealthCheck)
	}
}

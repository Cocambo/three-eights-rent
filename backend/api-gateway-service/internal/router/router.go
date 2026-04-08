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

	userProxy := proxy.NewReverseProxy(userServiceURL, "/api/v1/users", "/api/v1")
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

	protectedUsers := engine.Group("/api/v1/users")
	protectedUsers.Use(stripIdentityHeaders, jwtMiddleware)
	{
		protectedUsers.GET("/profile", userProxy)
		protectedUsers.PUT("/profile", userProxy)
		protectedUsers.POST("/driver-license", userProxy)
	}
}

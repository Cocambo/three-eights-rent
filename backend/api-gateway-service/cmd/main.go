package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"api-gateway-service/internal/config"
	"api-gateway-service/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	gin.SetMode(cfg.GinMode)

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	router.Register(engine, cfg)

	if err := engine.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("failed to run gateway: %v", err)
	}
}

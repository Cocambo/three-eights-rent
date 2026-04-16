package main

import (
	"context"
	"log"

	"car-service/config"
	"car-service/internal/handler"
	"car-service/internal/middleware"
	"car-service/internal/repository"
	"car-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	gin.SetMode(cfg.Server.GinMode)

	db, err := config.NewPostgresDB(ctx, cfg)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}

	minioClient, err := config.NewMinIOClient(cfg)
	if err != nil {
		log.Fatalf("connect minio: %v", err)
	}

	carRepository := repository.NewCarRepository(db)
	carService := service.NewCarService(carRepository)
	carHandler := handler.NewCarHandler(carService)
	jwtMiddleware := middleware.NewJWTMiddleware(cfg.JWT.AccessSecret)

	router := gin.New()
	middleware.Setup(router)
	handler.RegisterRoutes(router, handler.Dependencies{
		CarHandler:     carHandler,
		JWTMiddleware:  jwtMiddleware,
		MinIOAvailable: minioClient != nil,
	})

	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("run server: %v", err)
	}
}

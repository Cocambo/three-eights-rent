package main

import (
	"context"
	"fmt"
	"log"

	"car-service/config"
	"car-service/database"
	"car-service/internal/handler"
	"car-service/internal/middleware"
	"car-service/internal/repository"
	"car-service/internal/service"
	"car-service/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	gin.SetMode(cfg.Server.GinMode)

	db, err := database.Open(ctx, cfg.DB)
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}
	defer func() {
		if err := database.Close(db); err != nil {
			log.Printf("close postgres: %v", err)
		}
	}()

	minioClient, err := config.NewMinIOClient(cfg)
	if err != nil {
		return fmt.Errorf("connect minio: %w", err)
	}

	minioPublicClient, err := config.NewMinIOPublicClient(cfg)
	if err != nil {
		return fmt.Errorf("connect minio public endpoint: %w", err)
	}

	imageStorage := storage.NewMinIOImageStorageService(minioClient, minioPublicClient, cfg.MinIO)
	carRepository := repository.NewCarRepository(db)
	bookingRepository := repository.NewBookingRepository(db)
	favoriteRepository := repository.NewFavoriteRepository(db)
	recommendationRepository := repository.NewRecommendationRepository(db)
	carService := service.NewCarService(carRepository, imageStorage)
	bookingService := service.NewBookingService(bookingRepository, carRepository, imageStorage)
	favoriteService := service.NewFavoriteService(favoriteRepository, carRepository, imageStorage)
	recommendationCalculator := service.NewRecommendationCalculator()
	recommendationService := service.NewRecommendationService(
		recommendationRepository,
		carRepository,
		imageStorage,
		recommendationCalculator,
	)
	carHandler := handler.NewCarHandler(carService)
	bookingHandler := handler.NewBookingHandler(bookingService)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService)
	recommendationHandler := handler.NewRecommendationHandler(recommendationService)
	jwtMiddleware := middleware.NewJWTMiddleware(cfg.JWT.AccessSecret)

	router := gin.New()
	middleware.Setup(router)
	handler.RegisterRoutes(router, handler.Dependencies{
		CarHandler:            carHandler,
		BookingHandler:        bookingHandler,
		FavoriteHandler:       favoriteHandler,
		RecommendationHandler: recommendationHandler,
		JWTMiddleware:         jwtMiddleware,
		MinIOAvailable:        imageStorage != nil,
	})

	if err := router.Run(":" + cfg.Server.Port); err != nil {
		return fmt.Errorf("run server: %w", err)
	}

	return nil
}

package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"user-service/internal/config"
	"user-service/internal/handler"
	"user-service/internal/middleware"
	"user-service/internal/repository"
	"user-service/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := config.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Manual dependency injection.
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	router := gin.New()
	middleware.Setup(router)
	handler.RegisterRoutes(router, userHandler)

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

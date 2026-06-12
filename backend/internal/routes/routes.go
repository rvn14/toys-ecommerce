package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rvn14/toys-ecommerce/backend/internal/config"
	"github.com/rvn14/toys-ecommerce/backend/internal/handlers"
	"github.com/rvn14/toys-ecommerce/backend/internal/middleware"
	"github.com/rvn14/toys-ecommerce/backend/internal/repositories"
	"github.com/rvn14/toys-ecommerce/backend/internal/services"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, cfg config.ApConfig) {
	api := router.Group("/api/v1")

	healthHandler := handlers.NewHealthHandler(db)
	api.GET("/health", healthHandler.HealthCheck)
	api.GET("/health/db", healthHandler.DBHealthCheck)

	userRepository := repositories.NewUserRepository(db)
	authService := services.NewAuthService(
		userRepository,
		cfg.JWTSecret,
		cfg.AccessTokenExpiresMinutes,
	)

	authHandler := handlers.NewAuthHandler(authService)

	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/signup", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.GET("/me", middleware.AuthMiddleware(cfg.JWTSecret), authHandler.Me)
	}
}

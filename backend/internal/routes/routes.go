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

func RegisterRoutes(
	router *gin.Engine,
	db *gorm.DB,
	cfg config.AppConfig,
	tokenBlacklist services.TokenBlacklist,
) {
	api := router.Group("/api/v1")

	healthHandler := handlers.NewHealthHandler(db)

	api.GET("/health", healthHandler.HealthCheck)
	api.GET("/health/db", healthHandler.DBHealthCheck)

	userRepository := repositories.NewUserRepository(db)
	authService := services.NewAuthService(
		userRepository,
		cfg.JWTSecret,
		cfg.AccessTokenExpiresMinutes,
		tokenBlacklist,
	)
	authHandler := handlers.NewAuthHandler(authService)

	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)

		protectedAuthRoutes := authRoutes.Group("")
		protectedAuthRoutes.Use(middleware.AuthMiddleware(cfg.JWTSecret, tokenBlacklist))
		{
			protectedAuthRoutes.GET("/me", authHandler.Me)
			protectedAuthRoutes.POST("/logout", authHandler.Logout)
		}
	}

	adminRoutes := api.Group("/admin")
	adminRoutes.Use(
		middleware.AuthMiddleware(cfg.JWTSecret, tokenBlacklist),
		middleware.AdminOnly(),
	)
	{
		adminRoutes.GET("/ping", handlers.AdminPing)
	}
}

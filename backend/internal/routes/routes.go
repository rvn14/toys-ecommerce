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

	catalogRepository := repositories.NewCatalogRepository(db)
	catalogService := services.NewCatalogService(catalogRepository)
	catalogHandler := handlers.NewCatalogHandler(catalogService)

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

	api.GET("/categories", catalogHandler.GetCategories)
	api.GET("/brands", catalogHandler.GetBrands)
	api.GET("/products", catalogHandler.GetProducts)
	api.GET("/products/:id", catalogHandler.GetProductByID)

	adminRoutes := api.Group("/admin")
	adminRoutes.Use(
		middleware.AuthMiddleware(cfg.JWTSecret, tokenBlacklist),
		middleware.AdminOnly(),
	)
	{
		adminRoutes.GET("/ping", handlers.AdminPing)

		adminRoutes.POST("/categories", catalogHandler.CreateCategory)
		adminRoutes.PUT("/categories/:id", catalogHandler.UpdateCategory)
		adminRoutes.DELETE("/categories/:id", catalogHandler.DeleteCategory)

		adminRoutes.POST("/brands", catalogHandler.CreateBrand)
		adminRoutes.PUT("/brands/:id", catalogHandler.UpdateBrand)
		adminRoutes.DELETE("/brands/:id", catalogHandler.DeleteBrand)

		adminRoutes.POST("/products", catalogHandler.CreateProduct)
		adminRoutes.PUT("/products/:id", catalogHandler.UpdateProduct)
		adminRoutes.DELETE("/products/:id", catalogHandler.DeleteProduct)
	}
}

package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/rvn14/toys-ecommerce/backend/docs"
	"github.com/rvn14/toys-ecommerce/backend/internal/config"
	"github.com/rvn14/toys-ecommerce/backend/internal/database"
	"github.com/rvn14/toys-ecommerce/backend/internal/routes"
	"github.com/rvn14/toys-ecommerce/backend/internal/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Toys Store API
// @version 1.0
// @description REST API for a toys and collectibles e-commerce store.
// @description Built with Go, Gin, GORM, PostgreSQL, JWT authentication, and role-based authorization.
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	cfg := config.LoadConfig()

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Database migration failed:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tokenBlacklist := services.NewInMemoryTokenBlacklist()
	tokenBlacklist.StartCleanup(ctx, 10*time.Minute)

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.RegisterRoutes(router, db, cfg, tokenBlacklist)

	serverAddress := ":" + cfg.Port

	log.Printf("Starting Global Toys Store API in %s mode on port %s", cfg.AppEnv, cfg.Port)

	if err := router.Run(serverAddress); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

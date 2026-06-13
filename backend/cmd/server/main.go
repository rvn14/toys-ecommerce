package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rvn14/toys-ecommerce/backend/internal/config"
	"github.com/rvn14/toys-ecommerce/backend/internal/database"
	"github.com/rvn14/toys-ecommerce/backend/internal/routes"
	"github.com/rvn14/toys-ecommerce/backend/internal/services"

	"log"
)

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

	routes.RegisterRoutes(router, db, cfg, tokenBlacklist)

	serverAddress := ":" + cfg.Port

	log.Printf("Starting Global Toys Store API in %s mode on port %s", cfg.AppEnv, cfg.Port)

	if err := router.Run(serverAddress); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

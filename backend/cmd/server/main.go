package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rvn14/toys-ecommerce/backend/internal/config"
	"github.com/rvn14/toys-ecommerce/backend/internal/database"
	"github.com/rvn14/toys-ecommerce/backend/internal/routes"

	"log"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	router := gin.Default()

	routes.RegisterRoutes(router, db, cfg)

	log.Printf("Starting server on port %s in %s mode", cfg.Port, cfg.AppEnv)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

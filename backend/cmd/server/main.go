package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rvn14/toys-ecommerce/backend/internal/config"
	"github.com/rvn14/toys-ecommerce/backend/internal/routes"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	router := gin.Default()

	routes.RegisterRoutes(router)

	log.Printf("Starting server on port %s in %s mode", cfg.Port, cfg.AppEnv)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

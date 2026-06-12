package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rvn14/toys-ecommerce/backend/internal/handlers"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")

	api.GET("/health", handlers.HealthCheck)
}

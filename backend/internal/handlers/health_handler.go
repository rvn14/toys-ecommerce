package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type HealthHandler struct {
	DB *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{DB: db}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Toys Store API is Running!",
	})
}

func (h *HealthHandler) DBHealthCheck(c *gin.Context) {
	sqlDB, err := h.DB.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database connection error",
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database is not reachable",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Database is healthy",
	})
}

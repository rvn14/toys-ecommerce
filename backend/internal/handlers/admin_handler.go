package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminPing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Admin route is working",
	})
}

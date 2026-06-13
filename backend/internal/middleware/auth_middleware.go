package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rvn14/toys-ecommerce/backend/internal/services"
	"github.com/rvn14/toys-ecommerce/backend/internal/utils"
)

func AuthMiddleware(jwtSecret string, tokenBlacklist services.TokenBlacklist) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := utils.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		if tokenBlacklist.IsRevoked(claims.ID) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token has been logged out",
			})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("tokenID", claims.ID)
		c.Set("tokenExpiresAt", claims.ExpiresAt.Time)

		c.Next()
	}
}

package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	allowedRoleSet := make(map[string]struct{}, len(allowedRoles))

	for _, role := range allowedRoles {
		allowedRoleSet[role] = struct{}{}
	}

	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "User role is missing",
			})
			c.Abort()
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid user role",
			})
			c.Abort()
			return
		}

		if _, allowed := allowedRoleSet[role]; !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "You do not have permission to access this resource",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return RequireRole("admin")
}

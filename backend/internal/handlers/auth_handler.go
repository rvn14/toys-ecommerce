package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rvn14/toys-ecommerce/backend/internal/dto"
	"github.com/rvn14/toys-ecommerce/backend/internal/services"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request dto.RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
		return
	}

	response, err := h.authService.Register(request)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Email already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to register user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User registered successfully",
		"data":    response,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request dto.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
		return
	}

	response, err := h.authService.Login(request)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid email or password",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to login",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"data":    response,
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid user context",
		})
		return
	}

	response, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get current user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Current user loaded successfully",
		"data":    response,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	tokenIDValue, tokenIDExists := c.Get("tokenID")
	expiresAtValue, expiresAtExists := c.Get("tokenExpiresAt")

	if !tokenIDExists || !expiresAtExists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid token context",
		})
		return
	}

	tokenID, tokenIDOk := tokenIDValue.(string)
	expiresAt, expiresAtOk := expiresAtValue.(time.Time)

	if !tokenIDOk || !expiresAtOk {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid token context",
		})
		return
	}

	if err := h.authService.Logout(tokenID, expiresAt); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Failed to logout",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out successfully",
	})
}

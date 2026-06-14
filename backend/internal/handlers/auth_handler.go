package handlers

import (
	"errors"
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
		validationError(c, err)
		return
	}

	response, err := h.authService.Register(request)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			conflict(c, "Email already exists")
			return
		}

		internalError(c, "Failed to register user")
		return
	}

	created(c, "User registered successfully", response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request dto.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		validationError(c, err)
		return
	}

	response, err := h.authService.Login(request)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			unauthorized(c, "Invalid email or password")
			return
		}

		internalError(c, "Failed to login")
		return
	}

	ok(c, "Login successful", response)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		unauthorized(c, "Unauthorized")
		return
	}

	userID, userIDOK := userIDValue.(uint)
	if !userIDOK {
		unauthorized(c, "Invalid user context")
		return
	}

	response, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			notFound(c, "User not found")
			return
		}

		internalError(c, "Failed to get current user")
		return
	}

	ok(c, "Current user loaded successfully", response)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	tokenIDValue, tokenIDExists := c.Get("tokenID")
	expiresAtValue, expiresAtExists := c.Get("tokenExpiresAt")

	if !tokenIDExists || !expiresAtExists {
		unauthorized(c, "Invalid token context")
		return
	}

	tokenID, tokenIDOk := tokenIDValue.(string)
	expiresAt, expiresAtOk := expiresAtValue.(time.Time)

	if !tokenIDOk || !expiresAtOk {
		unauthorized(c, "Invalid token context")
		return
	}

	if err := h.authService.Logout(tokenID, expiresAt); err != nil {
		unauthorized(c, "Failed to logout")
		return
	}

	ok(c, "Logged out successfully", nil)
}

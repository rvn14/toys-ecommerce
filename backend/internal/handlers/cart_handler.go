package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rvn14/toys-ecommerce/backend/internal/dto"
	"github.com/rvn14/toys-ecommerce/backend/internal/services"
)

type CartHandler struct {
	cartService services.CartService
}

func NewCartHandler(cartService services.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID, valid := getAuthenticatedUserID(c)
	if !valid {
		return
	}

	response, err := h.cartService.GetCart(userID)
	if err != nil {
		handleCartError(c, err)
		return
	}

	ok(c, "Cart loaded successfully", response)
}

func (h *CartHandler) AddItem(c *gin.Context) {
	userID, valid := getAuthenticatedUserID(c)
	if !valid {
		return
	}

	var request dto.AddToCartRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		badRequest(c, "Invalid request body", err.Error())
		return
	}

	response, err := h.cartService.AddItem(userID, request)
	if err != nil {
		handleCartError(c, err)
		return
	}

	created(c, "Item added to cart successfully", response)
}

func (h *CartHandler) UpdateItemQuantity(c *gin.Context) {
	userID, valid := getAuthenticatedUserID(c)
	if !valid {
		return
	}

	itemID, valid := getCartItemIDParam(c)
	if !valid {
		return
	}

	var request dto.UpdateCartItemQuantityRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		badRequest(c, "Invalid request body", err.Error())
		return
	}

	response, err := h.cartService.UpdateItemQuantity(userID, itemID, request)
	if err != nil {
		handleCartError(c, err)
		return
	}

	ok(c, "Cart item quantity updated successfully", response)
}

func (h *CartHandler) RemoveItem(c *gin.Context) {
	userID, valid := getAuthenticatedUserID(c)
	if !valid {
		return
	}

	itemID, valid := getCartItemIDParam(c)
	if !valid {
		return
	}

	response, err := h.cartService.RemoveItem(userID, itemID)
	if err != nil {
		handleCartError(c, err)
		return
	}

	ok(c, "Cart item removed successfully", response)
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, valid := getAuthenticatedUserID(c)
	if !valid {
		return
	}

	response, err := h.cartService.ClearCart(userID)
	if err != nil {
		handleCartError(c, err)
		return
	}

	ok(c, "Cart cleared successfully", response)
}

func getAuthenticatedUserID(c *gin.Context) (uint, bool) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return 0, false
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid user context",
		})
		return 0, false
	}

	return userID, true
}

func getCartItemIDParam(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		badRequest(c, "Invalid cart item ID", nil)
		return 0, false
	}

	return uint(id), true
}

func handleCartError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrCartProductNotFound):
		notFound(c, "Product not found")

	case errors.Is(err, services.ErrCartItemNotFound):
		notFound(c, "Cart item not found")

	case errors.Is(err, services.ErrCartInsufficientStock):
		badRequest(c, "Requested quantity exceeds available stock", nil)

	case errors.Is(err, services.ErrCartProductNotActive):
		badRequest(c, "Product is not available for purchase", nil)

	default:
		internalError(c, "Cart operation failed")
	}
}

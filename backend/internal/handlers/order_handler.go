package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/rvn14/toys-ecommerce/backend/internal/dto"
	"github.com/rvn14/toys-ecommerce/backend/internal/services"
)

type OrderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, valid := getAuthenticatedUserID(c)
	if !valid {
		return
	}

	var request dto.CreateOrderRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		badRequest(c, "Invalid request body", err.Error())
		return
	}

	response, err := h.orderService.CreateOrder(userID, request)
	if err != nil {
		handleOrderError(c, err)
		return
	}

	created(c, "Order placed successfully", response)
}

func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	userID, valid := getAuthenticatedUserID(c)
	if !valid {
		return
	}

	response, err := h.orderService.GetMyOrders(userID)
	if err != nil {
		internalError(c, "Failed to load orders")
		return
	}

	ok(c, "Orders loaded successfully", response)
}

func (h *OrderHandler) GetMyOrderByID(c *gin.Context) {
	userID, valid := getAuthenticatedUserID(c)
	if !valid {
		return
	}

	orderID, valid := getIDParam(c)
	if !valid {
		return
	}

	response, err := h.orderService.GetMyOrderByID(userID, orderID)
	if err != nil {
		handleOrderError(c, err)
		return
	}

	ok(c, "Order loaded successfully", response)
}

func (h *OrderHandler) AdminGetOrders(c *gin.Context) {
	response, err := h.orderService.GetAllOrders()
	if err != nil {
		internalError(c, "Failed to load orders")
		return
	}

	ok(c, "Orders loaded successfully", response)
}

func (h *OrderHandler) AdminGetOrderByID(c *gin.Context) {
	orderID, valid := getIDParam(c)
	if !valid {
		return
	}

	response, err := h.orderService.GetOrderByID(orderID)
	if err != nil {
		handleOrderError(c, err)
		return
	}

	ok(c, "Order loaded successfully", response)
}

func (h *OrderHandler) AdminUpdateOrderStatus(c *gin.Context) {
	orderID, valid := getIDParam(c)
	if !valid {
		return
	}

	var request dto.UpdateOrderStatusRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		badRequest(c, "Invalid request body", err.Error())
		return
	}

	response, err := h.orderService.UpdateOrderStatus(orderID, request)
	if err != nil {
		handleOrderError(c, err)
		return
	}

	ok(c, "Order status updated successfully", response)
}

func handleOrderError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrOrderCartEmpty):
		badRequest(c, "Cart is empty", nil)

	case errors.Is(err, services.ErrOrderProductNotFound):
		notFound(c, "Product not found")

	case errors.Is(err, services.ErrOrderInsufficientStock):
		badRequest(c, "Some cart items exceed available stock", nil)

	case errors.Is(err, services.ErrOrderProductNotActive):
		badRequest(c, "Some cart items are no longer available", nil)

	case errors.Is(err, services.ErrOrderNotFound):
		notFound(c, "Order not found")

	case errors.Is(err, services.ErrInvalidPaymentMethod):
		badRequest(c, "Invalid payment method", gin.H{
			"allowedPaymentMethods": []string{
				"cash_on_delivery",
				"bank_transfer",
				"card",
			},
		})

	case errors.Is(err, services.ErrInvalidOrderStatus):
		badRequest(c, "Invalid order status", gin.H{
			"allowedOrderStatuses": []string{
				"pending",
				"processing",
				"shipped",
				"delivered",
				"cancelled",
			},
		})

	default:
		internalError(c, "Order operation failed")
	}
}

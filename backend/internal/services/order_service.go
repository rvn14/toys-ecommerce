package services

import (
	"encoding/json"
	"errors"

	"github.com/rvn14/toys-ecommerce/backend/internal/dto"
	"github.com/rvn14/toys-ecommerce/backend/internal/models"
	"github.com/rvn14/toys-ecommerce/backend/internal/repositories"
	"gorm.io/datatypes"
)

var (
	ErrOrderCartEmpty         = errors.New("cart is empty")
	ErrOrderProductNotFound   = errors.New("product not found")
	ErrOrderInsufficientStock = errors.New("insufficient stock")
	ErrOrderProductNotActive  = errors.New("product is not active")
	ErrOrderNotFound          = errors.New("order not found")
	ErrInvalidPaymentMethod   = errors.New("invalid payment method")
	ErrInvalidOrderStatus     = errors.New("invalid order status")
)

type OrderService interface {
	CreateOrder(userID uint, request dto.CreateOrderRequest) (*dto.OrderResponse, error)
	GetMyOrders(userID uint) ([]dto.OrderResponse, error)
	GetMyOrderByID(userID uint, orderID uint) (*dto.OrderResponse, error)
	GetAllOrders() ([]dto.OrderResponse, error)
	GetOrderByID(orderID uint) (*dto.OrderResponse, error)
	UpdateOrderStatus(orderID uint, request dto.UpdateOrderStatusRequest) (*dto.OrderResponse, error)
}

type orderService struct {
	repository repositories.OrderRepository
}

func NewOrderService(repository repositories.OrderRepository) OrderService {
	return &orderService{repository: repository}
}

func (s *orderService) CreateOrder(userID uint, request dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	paymentMethod := models.PaymentMethod(request.PaymentMethod)

	if !isValidPaymentMethod(paymentMethod) {
		return nil, ErrInvalidPaymentMethod
	}

	shippingAddress, err := buildShippingAddressJSON(request.ShippingAddress)
	if err != nil {
		return nil, err
	}

	order, err := s.repository.CreateOrderFromCart(userID, shippingAddress, paymentMethod)
	if err != nil {
		return nil, mapOrderRepositoryError(err)
	}

	response := mapOrderToResponse(*order)
	return &response, nil
}

func (s *orderService) GetMyOrders(userID uint) ([]dto.OrderResponse, error) {
	orders, err := s.repository.GetOrdersByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.OrderResponse, 0, len(orders))
	for _, order := range orders {
		responses = append(responses, mapOrderToResponse(order))
	}

	return responses, nil
}

func (s *orderService) GetMyOrderByID(userID uint, orderID uint) (*dto.OrderResponse, error) {
	order, err := s.repository.GetOrderByIDForUser(userID, orderID)
	if err != nil {
		return nil, mapOrderRepositoryError(err)
	}

	response := mapOrderToResponse(*order)
	return &response, nil
}

func (s *orderService) GetAllOrders() ([]dto.OrderResponse, error) {
	orders, err := s.repository.GetAllOrders()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.OrderResponse, 0, len(orders))
	for _, order := range orders {
		responses = append(responses, mapOrderToResponse(order))
	}

	return responses, nil
}

func (s *orderService) GetOrderByID(orderID uint) (*dto.OrderResponse, error) {
	order, err := s.repository.GetOrderByID(orderID)
	if err != nil {
		return nil, mapOrderRepositoryError(err)
	}

	response := mapOrderToResponse(*order)
	return &response, nil
}

func (s *orderService) UpdateOrderStatus(orderID uint, request dto.UpdateOrderStatusRequest) (*dto.OrderResponse, error) {
	status := models.OrderStatus(request.Status)

	if !isValidOrderStatus(status) {
		return nil, ErrInvalidOrderStatus
	}

	order, err := s.repository.UpdateOrderStatus(orderID, status)
	if err != nil {
		return nil, mapOrderRepositoryError(err)
	}

	response := mapOrderToResponse(*order)
	return &response, nil
}

func buildShippingAddressJSON(address dto.ShippingAddressRequest) (datatypes.JSON, error) {
	bytes, err := json.Marshal(address)
	if err != nil {
		return nil, err
	}

	return datatypes.JSON(bytes), nil
}

func isValidPaymentMethod(method models.PaymentMethod) bool {
	switch method {
	case models.PaymentMethodCashOnDelivery,
		models.PaymentMethodBankTransfer,
		models.PaymentMethodCard:
		return true
	default:
		return false
	}
}

func isValidOrderStatus(status models.OrderStatus) bool {
	switch status {
	case models.OrderStatusPending,
		models.OrderStatusProcessing,
		models.OrderStatusShipped,
		models.OrderStatusDelivered,
		models.OrderStatusCancelled:
		return true
	default:
		return false
	}
}

func mapOrderRepositoryError(err error) error {
	switch {
	case errors.Is(err, repositories.ErrOrderCartEmpty):
		return ErrOrderCartEmpty
	case errors.Is(err, repositories.ErrOrderProductNotFound):
		return ErrOrderProductNotFound
	case errors.Is(err, repositories.ErrOrderInsufficientStock):
		return ErrOrderInsufficientStock
	case errors.Is(err, repositories.ErrOrderProductNotActive):
		return ErrOrderProductNotActive
	case errors.Is(err, repositories.ErrOrderNotFound):
		return ErrOrderNotFound
	default:
		return err
	}
}

func mapOrderToResponse(order models.Order) dto.OrderResponse {
	items := make([]dto.OrderItemResponse, 0, len(order.Items))

	for _, item := range order.Items {
		items = append(items, dto.OrderItemResponse{
			ID:                      item.ID,
			ProductID:               item.ProductID,
			ProductNameSnapshot:     item.ProductNameSnapshot,
			ProductSlugSnapshot:     item.ProductSlugSnapshot,
			ProductImageURLSnapshot: item.ProductImageURLSnapshot,
			Quantity:                item.Quantity,
			UnitPriceCents:          item.UnitPriceCents,
			LineTotalCents:          item.LineTotalCents,
		})
	}

	return dto.OrderResponse{
		ID:              order.ID,
		OrderNumber:     order.OrderNumber,
		UserID:          order.UserID,
		ShippingAddress: json.RawMessage(order.ShippingAddress),
		SubtotalCents:   order.SubtotalCents,
		TotalCents:      order.TotalCents,
		PaymentMethod:   string(order.PaymentMethod),
		PaymentStatus:   string(order.PaymentStatus),
		OrderStatus:     string(order.OrderStatus),
		Items:           items,
		CreatedAt:       order.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

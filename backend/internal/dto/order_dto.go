package dto

import "encoding/json"

type ShippingAddressRequest struct {
	FullName     string `json:"fullName" binding:"required,min=2,max=100"`
	Phone        string `json:"phone" binding:"required,min=7,max=30"`
	AddressLine1 string `json:"addressLine1" binding:"required,min=3,max=255"`
	AddressLine2 string `json:"addressLine2"`
	City         string `json:"city" binding:"required,min=2,max=100"`
	State        string `json:"state"`
	PostalCode   string `json:"postalCode"`
	Country      string `json:"country" binding:"required,min=2,max=100"`
}

type CreateOrderRequest struct {
	ShippingAddress ShippingAddressRequest `json:"shippingAddress" binding:"required"`
	PaymentMethod   string                 `json:"paymentMethod" binding:"required"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type OrderItemResponse struct {
	ID                      uint   `json:"id"`
	ProductID               uint   `json:"productId"`
	ProductNameSnapshot     string `json:"productNameSnapshot"`
	ProductSlugSnapshot     string `json:"productSlugSnapshot"`
	ProductImageURLSnapshot string `json:"productImageUrlSnapshot"`
	Quantity                int    `json:"quantity"`
	UnitPriceCents          int64  `json:"unitPriceCents"`
	LineTotalCents          int64  `json:"lineTotalCents"`
}

type OrderResponse struct {
	ID              uint                `json:"id"`
	OrderNumber     string              `json:"orderNumber"`
	UserID          uint                `json:"userId"`
	ShippingAddress json.RawMessage     `json:"shippingAddress"`
	SubtotalCents   int64               `json:"subtotalCents"`
	TotalCents      int64               `json:"totalCents"`
	PaymentMethod   string              `json:"paymentMethod"`
	PaymentStatus   string              `json:"paymentStatus"`
	OrderStatus     string              `json:"orderStatus"`
	Items           []OrderItemResponse `json:"items"`
	CreatedAt       string              `json:"createdAt"`
}

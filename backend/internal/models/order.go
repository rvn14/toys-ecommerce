package models

import (
	"time"

	"gorm.io/datatypes"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

type PaymentMethod string

const (
	PaymentMethodCashOnDelivery PaymentMethod = "cash_on_delivery"
	PaymentMethodBankTransfer   PaymentMethod = "bank_transfer"
	PaymentMethodCard           PaymentMethod = "card"
)

type Order struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	OrderNumber     string         `gorm:"type:varchar(40);not null;uniqueIndex" json:"orderNumber"`
	UserID          uint           `gorm:"not null;index" json:"userId"`
	User            User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	ShippingAddress datatypes.JSON `gorm:"type:jsonb;not null" json:"shippingAddress"`

	SubtotalCents int64 `gorm:"not null" json:"subtotalCents"`
	TotalCents    int64 `gorm:"not null" json:"totalCents"`

	PaymentMethod PaymentMethod `gorm:"type:varchar(30);not null" json:"paymentMethod"`
	PaymentStatus PaymentStatus `gorm:"type:varchar(30);not null;default:'pending';index" json:"paymentStatus"`
	OrderStatus   OrderStatus   `gorm:"type:varchar(30);not null;default:'pending';index" json:"orderStatus"`

	Items []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;" json:"items"`

	CreatedAt time.Time `gorm:"index" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Order) TableName() string {
	return "orders"
}

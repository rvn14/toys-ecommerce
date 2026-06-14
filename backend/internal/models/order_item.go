package models

import "time"

type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `gorm:"not null;index" json:"orderId"`
	ProductID uint    `gorm:"not null;index" json:"productId"`
	Product   Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`

	ProductNameSnapshot     string `gorm:"type:varchar(150);not null" json:"productNameSnapshot"`
	ProductSlugSnapshot     string `gorm:"type:varchar(180);not null" json:"productSlugSnapshot"`
	ProductImageURLSnapshot string `gorm:"type:text" json:"productImageUrlSnapshot"`

	Quantity       int   `gorm:"not null" json:"quantity"`
	UnitPriceCents int64 `gorm:"not null" json:"unitPriceCents"`
	LineTotalCents int64 `gorm:"not null" json:"lineTotalCents"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (OrderItem) TableName() string {
	return "order_items"
}

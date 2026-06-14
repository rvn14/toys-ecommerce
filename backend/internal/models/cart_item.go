package models

import "time"

type CartItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CartID    uint      `gorm:"not null;index;uniqueIndex:idx_cart_product" json:"cartId"`
	Cart      Cart      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	ProductID uint      `gorm:"not null;index;uniqueIndex:idx_cart_product" json:"productId"`
	Product   Product   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"product"`
	Quantity  int       `gorm:"not null;default:1;check:quantity > 0" json:"quantity"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (CartItem) TableName() string {
	return "cart_items"
}

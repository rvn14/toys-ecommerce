package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ProductStatus string

const (
	ProductStatusDraft      ProductStatus = "draft"
	ProductStatusActive     ProductStatus = "active"
	ProductStatusOutOfStock ProductStatus = "out_of_stock"
	ProductStatusArchived   ProductStatus = "archived"
)

type Product struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	Name               string         `gorm:"type:varchar(150);not null;index" json:"name"`
	Slug               string         `gorm:"type:varchar(180);not null;uniqueIndex" json:"slug"`
	SKU                string         `gorm:"type:varchar(80);not null;uniqueIndex" json:"sku"`
	Description        string         `gorm:"type:text" json:"description"`
	PriceCents         int64          `gorm:"not null;index" json:"priceCents"`
	DiscountPriceCents *int64         `json:"discountPriceCents"`
	StockQuantity      int            `gorm:"not null;default:0" json:"stockQuantity"`
	Status             ProductStatus  `gorm:"type:varchar(30);not null;default:'draft';index" json:"status"`
	Specifications     datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"specifications"`

	CategoryID uint     `gorm:"not null;index" json:"categoryId"`
	Category   Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"category"`

	BrandID uint  `gorm:"not null;index" json:"brandId"`
	Brand   Brand `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"brand"`

	Images []ProductImage `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;" json:"images"`

	CreatedAt time.Time      `gorm:"index" json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Product) TableName() string {
	return "products"
}

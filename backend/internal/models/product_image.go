package models

import "time"

type ProductImage struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProductID   uint      `gorm:"not null;index" json:"productId"`
	URL         string    `gorm:"type:text;not null" json:"url"`
	StoragePath string    `gorm:"type:text" json:"storagePath"`
	AltText     string    `gorm:"type:varchar(255)" json:"altText"`
	SortOrder   int       `gorm:"not null;default:0" json:"sortOrder"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (ProductImage) TableName() string {
	return "product_images"
}

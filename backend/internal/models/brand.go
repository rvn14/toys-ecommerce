package models

import (
	"time"

	"gorm.io/gorm"
)

type Brand struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Slug        string         `gorm:"type:varchar(120);not null;uniqueIndex" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	Products    []Product      `gorm:"foreignKey:BrandID" json:"-"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Brand) TableName() string {
	return "brands"
}

package models

import "time"

type Cart struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"not null;uniqueIndex" json:"userId"`
	User      User       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Items     []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE;" json:"items"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

func (Cart) TableName() string {
	return "carts"
}

package models

import (
	"time"
)

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	Email        string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Role         UserRole  `gorm:"type:varchar(50);not null;default:'customer'" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

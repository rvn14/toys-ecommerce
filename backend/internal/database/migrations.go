package database

import (
	"log"

	"github.com/rvn14/toys-ecommerce/backend/internal/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Brand{},
		&models.Product{},
		&models.ProductImage{},
	)
	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")

	return nil
}

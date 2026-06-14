package repositories

import (
	"errors"

	"github.com/rvn14/toys-ecommerce/backend/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrCartProductNotFound   = errors.New("product not found")
	ErrCartItemNotFound      = errors.New("cart item not found")
	ErrCartInsufficientStock = errors.New("insufficient stock")
	ErrCartProductNotActive  = errors.New("product is not active")
)

type CartRepository interface {
	GetCartByUserID(userID uint) (*models.Cart, error)
	AddItem(userID uint, productID uint, quantity int) (*models.Cart, error)
	UpdateItemQuantity(userID uint, itemID uint, quantity int) (*models.Cart, error)
	RemoveItem(userID uint, itemID uint) (*models.Cart, error)
	ClearCart(userID uint) (*models.Cart, error)
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetCartByUserID(userID uint) (*models.Cart, error) {
	cart, err := r.getOrCreateCart(r.db, userID)
	if err != nil {
		return nil, err
	}

	return r.loadCart(cart.ID)
}

func (r *cartRepository) AddItem(userID uint, productID uint, quantity int) (*models.Cart, error) {
	var cartID uint

	err := r.db.Transaction(func(tx *gorm.DB) error {
		product, err := r.getActiveProductForUpdate(tx, productID)
		if err != nil {
			return err
		}

		cart, err := r.getOrCreateCart(tx, userID)
		if err != nil {
			return err
		}

		cartID = cart.ID

		var item models.CartItem
		err = tx.
			Where("cart_id = ? AND product_id = ?", cart.ID, productID).
			First(&item).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			if quantity > product.StockQuantity {
				return ErrCartInsufficientStock
			}

			item = models.CartItem{
				CartID:    cart.ID,
				ProductID: productID,
				Quantity:  quantity,
			}

			return tx.Create(&item).Error
		}

		newQuantity := item.Quantity + quantity

		if newQuantity > product.StockQuantity {
			return ErrCartInsufficientStock
		}

		item.Quantity = newQuantity

		return tx.Save(&item).Error
	})

	if err != nil {
		return nil, err
	}

	return r.loadCart(cartID)
}

func (r *cartRepository) UpdateItemQuantity(userID uint, itemID uint, quantity int) (*models.Cart, error) {
	var cartID uint

	err := r.db.Transaction(func(tx *gorm.DB) error {
		cart, err := r.getOrCreateCart(tx, userID)
		if err != nil {
			return err
		}

		cartID = cart.ID

		var item models.CartItem
		err = tx.
			Where("id = ? AND cart_id = ?", itemID, cart.ID).
			First(&item).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCartItemNotFound
		}

		if err != nil {
			return err
		}

		product, err := r.getActiveProductForUpdate(tx, item.ProductID)
		if err != nil {
			return err
		}

		if quantity > product.StockQuantity {
			return ErrCartInsufficientStock
		}

		item.Quantity = quantity

		return tx.Save(&item).Error
	})

	if err != nil {
		return nil, err
	}

	return r.loadCart(cartID)
}

func (r *cartRepository) RemoveItem(userID uint, itemID uint) (*models.Cart, error) {
	var cartID uint

	err := r.db.Transaction(func(tx *gorm.DB) error {
		cart, err := r.getOrCreateCart(tx, userID)
		if err != nil {
			return err
		}

		cartID = cart.ID

		result := tx.
			Where("id = ? AND cart_id = ?", itemID, cart.ID).
			Delete(&models.CartItem{})

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return ErrCartItemNotFound
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return r.loadCart(cartID)
}

func (r *cartRepository) ClearCart(userID uint) (*models.Cart, error) {
	var cartID uint

	err := r.db.Transaction(func(tx *gorm.DB) error {
		cart, err := r.getOrCreateCart(tx, userID)
		if err != nil {
			return err
		}

		cartID = cart.ID

		return tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error
	})

	if err != nil {
		return nil, err
	}

	return r.loadCart(cartID)
}

func (r *cartRepository) getOrCreateCart(db *gorm.DB, userID uint) (*models.Cart, error) {
	cart := models.Cart{
		UserID: userID,
	}

	err := db.FirstOrCreate(&cart, models.Cart{UserID: userID}).Error
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepository) getActiveProductForUpdate(db *gorm.DB, productID uint) (*models.Product, error) {
	var product models.Product

	err := db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&product, productID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrCartProductNotFound
	}

	if err != nil {
		return nil, err
	}

	if product.Status != models.ProductStatusActive {
		return nil, ErrCartProductNotActive
	}

	return &product, nil
}

func (r *cartRepository) loadCart(cartID uint) (*models.Cart, error) {
	var cart models.Cart

	err := r.db.
		Preload("Items.Product.Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		First(&cart, cartID).Error

	if err != nil {
		return nil, err
	}

	return &cart, nil
}

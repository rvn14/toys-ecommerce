package repositories

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rvn14/toys-ecommerce/backend/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrOrderCartEmpty         = errors.New("cart is empty")
	ErrOrderProductNotFound   = errors.New("product not found")
	ErrOrderInsufficientStock = errors.New("insufficient stock")
	ErrOrderProductNotActive  = errors.New("product is not active")
	ErrOrderNotFound          = errors.New("order not found")
)

type OrderRepository interface {
	CreateOrderFromCart(userID uint, shippingAddress datatypes.JSON, paymentMethod models.PaymentMethod) (*models.Order, error)
	GetOrdersByUserID(userID uint) ([]models.Order, error)
	GetOrderByIDForUser(userID uint, orderID uint) (*models.Order, error)
	GetAllOrders() ([]models.Order, error)
	GetOrderByID(orderID uint) (*models.Order, error)
	UpdateOrderStatus(orderID uint, status models.OrderStatus) (*models.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrderFromCart(
	userID uint,
	shippingAddress datatypes.JSON,
	paymentMethod models.PaymentMethod,
) (*models.Order, error) {
	var orderID uint

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var cart models.Cart

		err := tx.
			Preload("Items").
			Where("user_id = ?", userID).
			First(&cart).Error

		if errors.Is(err, gorm.ErrRecordNotFound) || len(cart.Items) == 0 {
			return ErrOrderCartEmpty
		}

		if err != nil {
			return err
		}

		order := models.Order{
			OrderNumber:     generateOrderNumber(),
			UserID:          userID,
			ShippingAddress: shippingAddress,
			PaymentMethod:   paymentMethod,
			PaymentStatus:   models.PaymentStatusPending,
			OrderStatus:     models.OrderStatusPending,
		}

		var orderItems []models.OrderItem
		var subtotalCents int64

		for _, cartItem := range cart.Items {
			var product models.Product

			err := tx.
				Clauses(clause.Locking{Strength: "UPDATE"}).
				Preload("Images", func(db *gorm.DB) *gorm.DB {
					return db.Order("sort_order ASC")
				}).
				First(&product, cartItem.ProductID).Error

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrOrderProductNotFound
			}

			if err != nil {
				return err
			}

			if product.Status != models.ProductStatusActive {
				return ErrOrderProductNotActive
			}

			if cartItem.Quantity > product.StockQuantity {
				return ErrOrderInsufficientStock
			}

			unitPrice := calculateOrderUnitPriceCents(product)
			lineTotal := unitPrice * int64(cartItem.Quantity)

			subtotalCents += lineTotal

			orderItems = append(orderItems, models.OrderItem{
				ProductID:               product.ID,
				ProductNameSnapshot:     product.Name,
				ProductSlugSnapshot:     product.Slug,
				ProductImageURLSnapshot: firstProductImageURLForOrder(product),
				Quantity:                cartItem.Quantity,
				UnitPriceCents:          unitPrice,
				LineTotalCents:          lineTotal,
			})

			product.StockQuantity -= cartItem.Quantity
			if product.StockQuantity == 0 {
				product.Status = models.ProductStatusOutOfStock
			}

			if err := tx.Save(&product).Error; err != nil {
				return err
			}
		}

		order.SubtotalCents = subtotalCents
		order.TotalCents = subtotalCents

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		for index := range orderItems {
			orderItems[index].OrderID = order.ID
		}

		if err := tx.Create(&orderItems).Error; err != nil {
			return err
		}

		if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}

		orderID = order.ID

		return nil
	})

	if err != nil {
		return nil, err
	}

	return r.GetOrderByID(orderID)
}

func (r *orderRepository) GetOrdersByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order

	err := r.db.
		Preload("Items").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, err
}

func (r *orderRepository) GetOrderByIDForUser(userID uint, orderID uint) (*models.Order, error) {
	var order models.Order

	err := r.db.
		Preload("Items").
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrOrderNotFound
	}

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order

	err := r.db.
		Preload("Items").
		Order("created_at DESC").
		Find(&orders).Error

	return orders, err
}

func (r *orderRepository) GetOrderByID(orderID uint) (*models.Order, error) {
	var order models.Order

	err := r.db.
		Preload("Items").
		First(&order, orderID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrOrderNotFound
	}

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) UpdateOrderStatus(orderID uint, status models.OrderStatus) (*models.Order, error) {
	order, err := r.GetOrderByID(orderID)
	if err != nil {
		return nil, err
	}

	order.OrderStatus = status

	if err := r.db.Save(order).Error; err != nil {
		return nil, err
	}

	return r.GetOrderByID(order.ID)
}

func generateOrderNumber() string {
	now := time.Now().UTC().Format("20060102150405")
	randomPart := uuid.NewString()[:8]

	return "ORD-" + now + "-" + randomPart
}

func calculateOrderUnitPriceCents(product models.Product) int64 {
	if product.DiscountPriceCents != nil && *product.DiscountPriceCents > 0 {
		return *product.DiscountPriceCents
	}

	return product.PriceCents
}

func firstProductImageURLForOrder(product models.Product) string {
	if len(product.Images) == 0 {
		return ""
	}

	return product.Images[0].URL
}

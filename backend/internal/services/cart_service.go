package services

import (
	"errors"

	"github.com/rvn14/toys-ecommerce/backend/internal/dto"
	"github.com/rvn14/toys-ecommerce/backend/internal/models"
	"github.com/rvn14/toys-ecommerce/backend/internal/repositories"
)

var (
	ErrCartProductNotFound   = errors.New("product not found")
	ErrCartItemNotFound      = errors.New("cart item not found")
	ErrCartInsufficientStock = errors.New("insufficient stock")
	ErrCartProductNotActive  = errors.New("product is not active")
)

type CartService interface {
	GetCart(userID uint) (*dto.CartResponse, error)
	AddItem(userID uint, request dto.AddToCartRequest) (*dto.CartResponse, error)
	UpdateItemQuantity(userID uint, itemID uint, request dto.UpdateCartItemQuantityRequest) (*dto.CartResponse, error)
	RemoveItem(userID uint, itemID uint) (*dto.CartResponse, error)
	ClearCart(userID uint) (*dto.CartResponse, error)
}

type cartService struct {
	repository repositories.CartRepository
}

func NewCartService(repository repositories.CartRepository) CartService {
	return &cartService{repository: repository}
}

func (s *cartService) GetCart(userID uint) (*dto.CartResponse, error) {
	cart, err := s.repository.GetCartByUserID(userID)
	if err != nil {
		return nil, mapCartRepositoryError(err)
	}

	response := mapCartToResponse(*cart)

	return &response, nil
}

func (s *cartService) AddItem(userID uint, request dto.AddToCartRequest) (*dto.CartResponse, error) {
	cart, err := s.repository.AddItem(userID, request.ProductID, request.Quantity)
	if err != nil {
		return nil, mapCartRepositoryError(err)
	}

	response := mapCartToResponse(*cart)

	return &response, nil
}

func (s *cartService) UpdateItemQuantity(userID uint, itemID uint, request dto.UpdateCartItemQuantityRequest) (*dto.CartResponse, error) {
	cart, err := s.repository.UpdateItemQuantity(userID, itemID, request.Quantity)
	if err != nil {
		return nil, mapCartRepositoryError(err)
	}

	response := mapCartToResponse(*cart)

	return &response, nil
}

func (s *cartService) RemoveItem(userID uint, itemID uint) (*dto.CartResponse, error) {
	cart, err := s.repository.RemoveItem(userID, itemID)
	if err != nil {
		return nil, mapCartRepositoryError(err)
	}

	response := mapCartToResponse(*cart)

	return &response, nil
}

func (s *cartService) ClearCart(userID uint) (*dto.CartResponse, error) {
	cart, err := s.repository.ClearCart(userID)
	if err != nil {
		return nil, mapCartRepositoryError(err)
	}

	response := mapCartToResponse(*cart)

	return &response, nil
}

func mapCartRepositoryError(err error) error {
	switch {
	case errors.Is(err, repositories.ErrCartProductNotFound):
		return ErrCartProductNotFound
	case errors.Is(err, repositories.ErrCartItemNotFound):
		return ErrCartItemNotFound
	case errors.Is(err, repositories.ErrCartInsufficientStock):
		return ErrCartInsufficientStock
	case errors.Is(err, repositories.ErrCartProductNotActive):
		return ErrCartProductNotActive
	default:
		return err
	}
}

func mapCartToResponse(cart models.Cart) dto.CartResponse {
	items := make([]dto.CartItemResponse, 0, len(cart.Items))

	var totalItems int
	var subtotalCents int64

	for _, item := range cart.Items {
		effectivePrice := getEffectivePriceCents(item.Product)
		lineTotal := effectivePrice * int64(item.Quantity)
		imageURL := getFirstProductImageURL(item.Product)

		totalItems += item.Quantity
		subtotalCents += lineTotal

		items = append(items, dto.CartItemResponse{
			ID:                  item.ID,
			ProductID:           item.ProductID,
			ProductName:         item.Product.Name,
			ProductSlug:         item.Product.Slug,
			ProductImageURL:     imageURL,
			PriceCents:          item.Product.PriceCents,
			DiscountPriceCents:  item.Product.DiscountPriceCents,
			EffectivePriceCents: effectivePrice,
			StockQuantity:       item.Product.StockQuantity,
			Quantity:            item.Quantity,
			LineTotalCents:      lineTotal,
		})
	}

	return dto.CartResponse{
		ID:            cart.ID,
		UserID:        cart.UserID,
		Items:         items,
		TotalItems:    totalItems,
		SubtotalCents: subtotalCents,
	}
}

func getEffectivePriceCents(product models.Product) int64 {
	if product.DiscountPriceCents != nil && *product.DiscountPriceCents > 0 {
		return *product.DiscountPriceCents
	}

	return product.PriceCents
}

func getFirstProductImageURL(product models.Product) string {
	if len(product.Images) == 0 {
		return ""
	}

	return product.Images[0].URL
}

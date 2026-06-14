package dto

type AddToCartRequest struct {
	ProductID uint `json:"productId" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1,max=99"`
}

type UpdateCartItemQuantityRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1,max=99"`
}

type CartItemResponse struct {
	ID                  uint   `json:"id"`
	ProductID           uint   `json:"productId"`
	ProductName         string `json:"productName"`
	ProductSlug         string `json:"productSlug"`
	ProductImageURL     string `json:"productImageUrl"`
	PriceCents          int64  `json:"priceCents"`
	DiscountPriceCents  *int64 `json:"discountPriceCents"`
	EffectivePriceCents int64  `json:"effectivePriceCents"`
	StockQuantity       int    `json:"stockQuantity"`
	Quantity            int    `json:"quantity"`
	LineTotalCents      int64  `json:"lineTotalCents"`
}

type CartResponse struct {
	ID            uint               `json:"id"`
	UserID        uint               `json:"userId"`
	Items         []CartItemResponse `json:"items"`
	TotalItems    int                `json:"totalItems"`
	SubtotalCents int64              `json:"subtotalCents"`
}

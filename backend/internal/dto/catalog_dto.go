package dto

import "encoding/json"

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
}

type CreateBrandRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
}

type UpdateBrandRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
}

type ProductImageInput struct {
	URL         string `json:"url" binding:"required"`
	StoragePath string `json:"storagePath"`
	AltText     string `json:"altText"`
	SortOrder   int    `json:"sortOrder"`
}

type CreateProductRequest struct {
	Name               string                 `json:"name" binding:"required,min=2,max=150"`
	SKU                string                 `json:"sku" binding:"required,min=2,max=80"`
	Description        string                 `json:"description"`
	PriceCents         int64                  `json:"priceCents" binding:"required,min=1"`
	DiscountPriceCents *int64                 `json:"discountPriceCents"`
	StockQuantity      int                    `json:"stockQuantity" binding:"gte=0"`
	Status             string                 `json:"status"`
	CategoryID         uint                   `json:"categoryId" binding:"required"`
	BrandID            uint                   `json:"brandId" binding:"required"`
	Specifications     map[string]interface{} `json:"specifications"`
	Images             []ProductImageInput    `json:"images"`
}

type UpdateProductRequest = CreateProductRequest

type CategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type BrandResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type ProductImageResponse struct {
	ID          uint   `json:"id"`
	URL         string `json:"url"`
	StoragePath string `json:"storagePath"`
	AltText     string `json:"altText"`
	SortOrder   int    `json:"sortOrder"`
}

type ProductResponse struct {
	ID                 uint                   `json:"id"`
	Name               string                 `json:"name"`
	Slug               string                 `json:"slug"`
	SKU                string                 `json:"sku"`
	Description        string                 `json:"description"`
	PriceCents         int64                  `json:"priceCents"`
	DiscountPriceCents *int64                 `json:"discountPriceCents"`
	StockQuantity      int                    `json:"stockQuantity"`
	Status             string                 `json:"status"`
	Category           CategoryResponse       `json:"category"`
	Brand              BrandResponse          `json:"brand"`
	Specifications     json.RawMessage        `json:"specifications"`
	Images             []ProductImageResponse `json:"images"`
}

type ProductQueryParams struct {
	Search        string `form:"search"`
	Category      string `form:"category"`
	Brand         string `form:"brand"`
	CategoryID    uint   `form:"categoryId"`
	BrandID       uint   `form:"brandId"`
	MinPriceCents int64  `form:"minPriceCents"`
	MaxPriceCents int64  `form:"maxPriceCents"`
	Sort          string `form:"sort"`
	Page          int    `form:"page"`
	Limit         int    `form:"limit"`
}

type PaginationMeta struct {
	Page        int   `json:"page"`
	Limit       int   `json:"limit"`
	TotalItems  int64 `json:"totalItems"`
	TotalPages  int   `json:"totalPages"`
	HasNext     bool  `json:"hasNext"`
	HasPrevious bool  `json:"hasPrevious"`
}

type ProductListResponse struct {
	Items      []ProductResponse `json:"items"`
	Pagination PaginationMeta    `json:"pagination"`
}

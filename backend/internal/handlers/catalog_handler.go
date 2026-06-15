package handlers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rvn14/toys-ecommerce/backend/internal/dto"
	"github.com/rvn14/toys-ecommerce/backend/internal/services"
	"gorm.io/gorm"
)

type CatalogHandler struct {
	catalogService services.CatalogService
}

func NewCatalogHandler(catalogService services.CatalogService) *CatalogHandler {
	return &CatalogHandler{catalogService: catalogService}
}

func (h *CatalogHandler) CreateCategory(c *gin.Context) {
	var request dto.CreateCategoryRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		validationError(c, err)
		return
	}

	response, err := h.catalogService.CreateCategory(request)
	if err != nil {
		internalError(c, "Failed to create category")
		return
	}

	created(c, "Category created successfully", response)
}

func (h *CatalogHandler) GetCategories(c *gin.Context) {
	response, err := h.catalogService.GetCategories()
	if err != nil {
		internalError(c, "Failed to load categories")
		return
	}

	ok(c, "Categories loaded successfully", response)
}

func (h *CatalogHandler) UpdateCategory(c *gin.Context) {
	id, valid := getIDParam(c)
	if !valid {
		return
	}

	var request dto.UpdateCategoryRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		validationError(c, err)
		return
	}

	response, err := h.catalogService.UpdateCategory(id, request)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFound(c, "Category not found")
			return
		}

		internalError(c, "Failed to update category")
		return
	}

	ok(c, "Category updated successfully", response)
}

func (h *CatalogHandler) DeleteCategory(c *gin.Context) {
	id, valid := getIDParam(c)
	if !valid {
		return
	}

	if err := h.catalogService.DeleteCategory(id); err != nil {
		internalError(c, "Failed to delete category")
		return
	}

	ok(c, "Category deleted successfully", nil)
}

func (h *CatalogHandler) CreateBrand(c *gin.Context) {
	var request dto.CreateBrandRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		validationError(c, err)
		return
	}

	response, err := h.catalogService.CreateBrand(request)
	if err != nil {
		internalError(c, "Failed to create brand")
		return
	}

	created(c, "Brand created successfully", response)
}

func (h *CatalogHandler) GetBrands(c *gin.Context) {
	response, err := h.catalogService.GetBrands()
	if err != nil {
		internalError(c, "Failed to load brands")
		return
	}

	ok(c, "Brands loaded successfully", response)
}

func (h *CatalogHandler) UpdateBrand(c *gin.Context) {
	id, valid := getIDParam(c)
	if !valid {
		return
	}

	var request dto.UpdateBrandRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		validationError(c, err)
		return
	}

	response, err := h.catalogService.UpdateBrand(id, request)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFound(c, "Brand not found")
			return
		}

		internalError(c, "Failed to update brand")
		return
	}

	ok(c, "Brand updated successfully", response)
}

func (h *CatalogHandler) DeleteBrand(c *gin.Context) {
	id, valid := getIDParam(c)
	if !valid {
		return
	}

	if err := h.catalogService.DeleteBrand(id); err != nil {
		internalError(c, "Failed to delete brand")
		return
	}

	ok(c, "Brand deleted successfully", nil)
}

func (h *CatalogHandler) CreateProduct(c *gin.Context) {
	var request dto.CreateProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		validationError(c, err)
		return
	}

	response, err := h.catalogService.CreateProduct(request)
	if err != nil {
		if errors.Is(err, services.ErrInvalidProductStatus) {
			badRequest(c, "Invalid product status", nil)
			return
		}

		if errors.Is(err, services.ErrInvalidDiscountPrice) {
			badRequest(c, "Discount price must be lower than price", nil)
			return
		}

		internalError(c, "Failed to create product")
		return
	}

	created(c, "Product created successfully", response)
}

func (h *CatalogHandler) GetProducts(c *gin.Context) {
	h.getProducts(c, false)
}

func (h *CatalogHandler) GetAdminProducts(c *gin.Context) {
	h.getProducts(c, true)
}

func (h *CatalogHandler) getProducts(c *gin.Context, includeInactive bool) {
	var query dto.ProductQueryParams

	if err := c.ShouldBindQuery(&query); err != nil {
		validationError(c, err)
		return
	}

	var response *dto.ProductListResponse
	var err error
	if includeInactive {
		response, err = h.catalogService.GetAdminProducts(query)
	} else {
		response, err = h.catalogService.GetProducts(query)
	}
	if err != nil {
		if errors.Is(err, services.ErrInvalidSortOption) {
			badRequest(c, "Invalid sort option", gin.H{
				"allowedSortOptions": []string{
					"newest",
					"price_asc",
					"price_desc",
					"name_asc",
					"name_desc",
				},
			})
			return
		}

		internalError(c, "Failed to load products")
		return
	}

	ok(c, "Products loaded successfully", response)
}

func (h *CatalogHandler) GetProductByID(c *gin.Context) {
	id, valid := getIDParam(c)
	if !valid {
		return
	}

	response, err := h.catalogService.GetProductByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFound(c, "Product not found")
			return
		}

		internalError(c, "Failed to load product")
		return
	}

	ok(c, "Product loaded successfully", response)
}

func (h *CatalogHandler) UpdateProduct(c *gin.Context) {
	id, valid := getIDParam(c)
	if !valid {
		return
	}

	var request dto.UpdateProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		badRequest(c, "Invalid request body", err.Error())
		return
	}

	response, err := h.catalogService.UpdateProduct(id, request)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFound(c, "Product not found")
			return
		}

		if errors.Is(err, services.ErrInvalidProductStatus) {
			badRequest(c, "Invalid product status", nil)
			return
		}

		if errors.Is(err, services.ErrInvalidDiscountPrice) {
			badRequest(c, "Discount price must be lower than price", nil)
			return
		}

		internalError(c, "Failed to update product")
		return
	}

	ok(c, "Product updated successfully", response)
}

func (h *CatalogHandler) DeleteProduct(c *gin.Context) {
	id, valid := getIDParam(c)
	if !valid {
		return
	}

	if err := h.catalogService.DeleteProduct(id); err != nil {
		internalError(c, "Failed to delete product")
		return
	}

	ok(c, "Product deleted successfully", nil)
}

func getIDParam(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		badRequest(c, "Invalid ID parameter", nil)
		return 0, false
	}

	return uint(id), true
}

package services

import (
	"encoding/json"
	"errors"

	"github.com/rvn14/toys-ecommerce/backend/internal/dto"
	"github.com/rvn14/toys-ecommerce/backend/internal/models"
	"github.com/rvn14/toys-ecommerce/backend/internal/repositories"
	"github.com/rvn14/toys-ecommerce/backend/internal/utils"
	"gorm.io/datatypes"
)

var (
	ErrInvalidProductStatus = errors.New("invalid product status")
	ErrInvalidSortOption    = errors.New("invalid sort option")
	ErrInvalidDiscountPrice = errors.New("discount price must be lower than price")
)

type CatalogService interface {
	CreateCategory(request dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetCategories() ([]dto.CategoryResponse, error)
	UpdateCategory(id uint, request dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	DeleteCategory(id uint) error

	CreateBrand(request dto.CreateBrandRequest) (*dto.BrandResponse, error)
	GetBrands() ([]dto.BrandResponse, error)
	UpdateBrand(id uint, request dto.UpdateBrandRequest) (*dto.BrandResponse, error)
	DeleteBrand(id uint) error

	CreateProduct(request dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetProducts(query dto.ProductQueryParams) (*dto.ProductListResponse, error)
	GetProductByID(id uint) (*dto.ProductResponse, error)
	UpdateProduct(id uint, request dto.UpdateProductRequest) (*dto.ProductResponse, error)
	DeleteProduct(id uint) error
}

type catalogService struct {
	repository repositories.CatalogRepository
}

func NewCatalogService(repository repositories.CatalogRepository) CatalogService {
	return &catalogService{repository: repository}
}

func (s *catalogService) CreateCategory(request dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category := &models.Category{
		Name:        request.Name,
		Slug:        utils.Slugify(request.Name),
		Description: request.Description,
	}

	if err := s.repository.CreateCategory(category); err != nil {
		return nil, err
	}

	response := mapCategoryToResponse(*category)
	return &response, nil
}

func (s *catalogService) GetCategories() ([]dto.CategoryResponse, error) {
	categories, err := s.repository.GetCategories()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		responses = append(responses, mapCategoryToResponse(category))
	}

	return responses, nil
}

func (s *catalogService) UpdateCategory(id uint, request dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.repository.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}

	category.Name = request.Name
	category.Slug = utils.Slugify(request.Name)
	category.Description = request.Description

	if err := s.repository.UpdateCategory(category); err != nil {
		return nil, err
	}

	response := mapCategoryToResponse(*category)
	return &response, nil
}

func (s *catalogService) DeleteCategory(id uint) error {
	return s.repository.DeleteCategory(id)
}

func (s *catalogService) CreateBrand(request dto.CreateBrandRequest) (*dto.BrandResponse, error) {
	brand := &models.Brand{
		Name:        request.Name,
		Slug:        utils.Slugify(request.Name),
		Description: request.Description,
	}

	if err := s.repository.CreateBrand(brand); err != nil {
		return nil, err
	}

	response := mapBrandToResponse(*brand)
	return &response, nil
}

func (s *catalogService) GetBrands() ([]dto.BrandResponse, error) {
	brands, err := s.repository.GetBrands()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.BrandResponse, 0, len(brands))
	for _, brand := range brands {
		responses = append(responses, mapBrandToResponse(brand))
	}

	return responses, nil
}

func (s *catalogService) UpdateBrand(id uint, request dto.UpdateBrandRequest) (*dto.BrandResponse, error) {
	brand, err := s.repository.GetBrandByID(id)
	if err != nil {
		return nil, err
	}

	brand.Name = request.Name
	brand.Slug = utils.Slugify(request.Name)
	brand.Description = request.Description

	if err := s.repository.UpdateBrand(brand); err != nil {
		return nil, err
	}

	response := mapBrandToResponse(*brand)
	return &response, nil
}

func (s *catalogService) DeleteBrand(id uint) error {
	return s.repository.DeleteBrand(id)
}

func (s *catalogService) CreateProduct(request dto.CreateProductRequest) (*dto.ProductResponse, error) {
	status := normalizeProductStatus(request.Status)

	if !isValidProductStatus(status) {
		return nil, ErrInvalidProductStatus
	}

	if err := validateProductPricing(request.PriceCents, request.DiscountPriceCents); err != nil {
		return nil, err
	}

	specifications, err := buildSpecificationsJSON(request.Specifications)
	if err != nil {
		return nil, err
	}

	product := &models.Product{
		Name:               request.Name,
		Slug:               utils.Slugify(request.Name),
		SKU:                request.SKU,
		Description:        request.Description,
		PriceCents:         request.PriceCents,
		DiscountPriceCents: request.DiscountPriceCents,
		StockQuantity:      request.StockQuantity,
		Status:             status,
		CategoryID:         request.CategoryID,
		BrandID:            request.BrandID,
		Specifications:     specifications,
		Images:             mapImageInputsToModels(request.Images),
	}

	if err := s.repository.CreateProduct(product); err != nil {
		return nil, err
	}

	savedProduct, err := s.repository.GetProductByID(product.ID)
	if err != nil {
		return nil, err
	}

	response := mapProductToResponse(*savedProduct)
	return &response, nil
}

func (s *catalogService) GetProducts(query dto.ProductQueryParams) (*dto.ProductListResponse, error) {
	page := normalizePage(query.Page)
	limit := normalizeLimit(query.Limit)
	sort := normalizeSort(query.Sort)

	if !isValidSortOption(sort) {
		return nil, ErrInvalidSortOption
	}

	options := repositories.ProductQueryOptions{
		Search:        query.Search,
		CategorySlug:  query.Category,
		BrandSlug:     query.Brand,
		CategoryID:    query.CategoryID,
		BrandID:       query.BrandID,
		MinPriceCents: query.MinPriceCents,
		MaxPriceCents: query.MaxPriceCents,
		Status:        string(models.ProductStatusActive),
		Sort:          sort,
		Page:          page,
		Limit:         limit,
	}

	result, err := s.repository.GetProducts(options)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ProductResponse, 0, len(result.Products))
	for _, product := range result.Products {
		items = append(items, mapProductToResponse(product))
	}

	totalPages := calculateTotalPages(result.Total, limit)

	response := &dto.ProductListResponse{
		Items: items,
		Pagination: dto.PaginationMeta{
			Page:        page,
			Limit:       limit,
			TotalItems:  result.Total,
			TotalPages:  totalPages,
			HasNext:     page < totalPages,
			HasPrevious: page > 1,
		},
	}

	return response, nil
}

func (s *catalogService) GetProductByID(id uint) (*dto.ProductResponse, error) {
	product, err := s.repository.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	response := mapProductToResponse(*product)
	return &response, nil
}

func (s *catalogService) UpdateProduct(id uint, request dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	status := normalizeProductStatus(request.Status)

	if !isValidProductStatus(status) {
		return nil, ErrInvalidProductStatus
	}

	if err := validateProductPricing(request.PriceCents, request.DiscountPriceCents); err != nil {
		return nil, err
	}

	specifications, err := buildSpecificationsJSON(request.Specifications)
	if err != nil {
		return nil, err
	}

	product, err := s.repository.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	product.Name = request.Name
	product.Slug = utils.Slugify(request.Name)
	product.SKU = request.SKU
	product.Description = request.Description
	product.PriceCents = request.PriceCents
	product.DiscountPriceCents = request.DiscountPriceCents
	product.StockQuantity = request.StockQuantity
	product.Status = status
	product.CategoryID = request.CategoryID
	product.BrandID = request.BrandID
	product.Specifications = specifications

	if err := s.repository.UpdateProduct(product); err != nil {
		return nil, err
	}

	if err := s.repository.ReplaceProductImages(product.ID, mapImageInputsToModels(request.Images)); err != nil {
		return nil, err
	}

	updatedProduct, err := s.repository.GetProductByID(product.ID)
	if err != nil {
		return nil, err
	}

	response := mapProductToResponse(*updatedProduct)
	return &response, nil
}

func (s *catalogService) DeleteProduct(id uint) error {
	return s.repository.DeleteProduct(id)
}

func normalizeProductStatus(status string) models.ProductStatus {
	if status == "" {
		return models.ProductStatusDraft
	}

	return models.ProductStatus(status)
}

func isValidProductStatus(status models.ProductStatus) bool {
	switch status {
	case models.ProductStatusDraft,
		models.ProductStatusActive,
		models.ProductStatusOutOfStock,
		models.ProductStatusArchived:
		return true
	default:
		return false
	}
}

func validateProductPricing(priceCents int64, discountPriceCents *int64) error {
	if discountPriceCents == nil {
		return nil
	}

	if *discountPriceCents <= 0 {
		return ErrInvalidDiscountPrice
	}

	if *discountPriceCents >= priceCents {
		return ErrInvalidDiscountPrice
	}

	return nil
}

func buildSpecificationsJSON(specifications map[string]interface{}) (datatypes.JSON, error) {
	if specifications == nil {
		specifications = map[string]interface{}{}
	}

	bytes, err := json.Marshal(specifications)
	if err != nil {
		return nil, err
	}

	return datatypes.JSON(bytes), nil
}

func mapImageInputsToModels(inputs []dto.ProductImageInput) []models.ProductImage {
	images := make([]models.ProductImage, 0, len(inputs))

	for _, input := range inputs {
		images = append(images, models.ProductImage{
			URL:         input.URL,
			StoragePath: input.StoragePath,
			AltText:     input.AltText,
			SortOrder:   input.SortOrder,
		})
	}

	return images
}

func mapCategoryToResponse(category models.Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
	}
}

func mapBrandToResponse(brand models.Brand) dto.BrandResponse {
	return dto.BrandResponse{
		ID:          brand.ID,
		Name:        brand.Name,
		Slug:        brand.Slug,
		Description: brand.Description,
	}
}

func mapProductToResponse(product models.Product) dto.ProductResponse {
	images := make([]dto.ProductImageResponse, 0, len(product.Images))

	for _, image := range product.Images {
		images = append(images, dto.ProductImageResponse{
			ID:          image.ID,
			URL:         image.URL,
			StoragePath: image.StoragePath,
			AltText:     image.AltText,
			SortOrder:   image.SortOrder,
		})
	}

	return dto.ProductResponse{
		ID:                 product.ID,
		Name:               product.Name,
		Slug:               product.Slug,
		SKU:                product.SKU,
		Description:        product.Description,
		PriceCents:         product.PriceCents,
		DiscountPriceCents: product.DiscountPriceCents,
		StockQuantity:      product.StockQuantity,
		Status:             string(product.Status),
		Category:           mapCategoryToResponse(product.Category),
		Brand:              mapBrandToResponse(product.Brand),
		Specifications:     json.RawMessage(product.Specifications),
		Images:             images,
	}
}

func normalizePage(page int) int {
	if page < 1 {
		return 1
	}

	return page
}

func normalizeLimit(limit int) int {
	if limit < 1 {
		return 12
	}

	if limit > 100 {
		return 100
	}

	return limit
}

func normalizeSort(sort string) string {
	if sort == "" {
		return "newest"
	}

	return sort
}

func isValidSortOption(sort string) bool {
	switch sort {
	case "newest", "price_asc", "price_desc", "name_asc", "name_desc":
		return true
	default:
		return false
	}
}

func calculateTotalPages(totalItems int64, limit int) int {
	if totalItems == 0 {
		return 0
	}

	return int((totalItems + int64(limit) - 1) / int64(limit))
}

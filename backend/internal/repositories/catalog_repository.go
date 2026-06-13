package repositories

import (
	"strings"

	"github.com/rvn14/toys-ecommerce/backend/internal/models"
	"gorm.io/gorm"
)

type ProductQueryOptions struct {
	Search        string
	CategorySlug  string
	BrandSlug     string
	CategoryID    uint
	BrandID       uint
	MinPriceCents int64
	MaxPriceCents int64
	Status        string
	Sort          string
	Page          int
	Limit         int
}

type ProductListResult struct {
	Products []models.Product
	Total    int64
	Page     int
	Limit    int
}

type CatalogRepository interface {
	CreateCategory(category *models.Category) error
	GetCategories() ([]models.Category, error)
	GetCategoryByID(id uint) (*models.Category, error)
	UpdateCategory(category *models.Category) error
	DeleteCategory(id uint) error

	CreateBrand(brand *models.Brand) error
	GetBrands() ([]models.Brand, error)
	GetBrandByID(id uint) (*models.Brand, error)
	UpdateBrand(brand *models.Brand) error
	DeleteBrand(id uint) error

	CreateProduct(product *models.Product) error
	GetProducts(options ProductQueryOptions) (*ProductListResult, error)
	GetProductByID(id uint) (*models.Product, error)
	UpdateProduct(product *models.Product) error
	DeleteProduct(id uint) error
	ReplaceProductImages(productID uint, images []models.ProductImage) error
}

type catalogRepository struct {
	db *gorm.DB
}

func NewCatalogRepository(db *gorm.DB) CatalogRepository {
	return &catalogRepository{db: db}
}

func (r *catalogRepository) CreateCategory(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *catalogRepository) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Order("name ASC").Find(&categories).Error
	return categories, err
}

func (r *catalogRepository) GetCategoryByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	return &category, err
}

func (r *catalogRepository) UpdateCategory(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *catalogRepository) DeleteCategory(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

func (r *catalogRepository) CreateBrand(brand *models.Brand) error {
	return r.db.Create(brand).Error
}

func (r *catalogRepository) GetBrands() ([]models.Brand, error) {
	var brands []models.Brand
	err := r.db.Order("name ASC").Find(&brands).Error
	return brands, err
}

func (r *catalogRepository) GetBrandByID(id uint) (*models.Brand, error) {
	var brand models.Brand
	err := r.db.First(&brand, id).Error
	return &brand, err
}

func (r *catalogRepository) UpdateBrand(brand *models.Brand) error {
	return r.db.Save(brand).Error
}

func (r *catalogRepository) DeleteBrand(id uint) error {
	return r.db.Delete(&models.Brand{}, id).Error
}

func (r *catalogRepository) CreateProduct(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *catalogRepository) GetProducts(options ProductQueryOptions) (*ProductListResult, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{})

	if options.Status != "" {
		query = query.Where("products.status = ?", options.Status)
	}

	if options.Search != "" {
		searchPattern := "%" + strings.ToLower(strings.TrimSpace(options.Search)) + "%"

		query = query.Where(
			"LOWER(products.name) LIKE ? OR LOWER(products.description) LIKE ? OR LOWER(products.sku) LIKE ?",
			searchPattern,
			searchPattern,
			searchPattern,
		)
	}

	if options.CategoryID != 0 {
		query = query.Where("products.category_id = ?", options.CategoryID)
	}

	if options.BrandID != 0 {
		query = query.Where("products.brand_id = ?", options.BrandID)
	}

	if options.CategorySlug != "" {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.slug = ?", options.CategorySlug)
	}

	if options.BrandSlug != "" {
		query = query.Joins("JOIN brands ON brands.id = products.brand_id").
			Where("brands.slug = ?", options.BrandSlug)
	}

	if options.MinPriceCents > 0 {
		query = query.Where("products.price_cents >= ?", options.MinPriceCents)
	}

	if options.MaxPriceCents > 0 {
		query = query.Where("products.price_cents <= ?", options.MaxPriceCents)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	orderBy := getProductOrderBy(options.Sort)
	offset := (options.Page - 1) * options.Limit

	err := query.
		Preload("Category").
		Preload("Brand").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Order(orderBy).
		Limit(options.Limit).
		Offset(offset).
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	return &ProductListResult{
		Products: products,
		Total:    total,
		Page:     options.Page,
		Limit:    options.Limit,
	}, nil
}

func getProductOrderBy(sort string) string {
	switch sort {
	case "price_asc":
		return "products.price_cents ASC"
	case "price_desc":
		return "products.price_cents DESC"
	case "name_asc":
		return "products.name ASC"
	case "name_desc":
		return "products.name DESC"
	case "newest":
		return "products.created_at DESC"
	default:
		return "products.created_at DESC"
	}
}

func (r *catalogRepository) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product

	err := r.db.
		Preload("Category").
		Preload("Brand").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		First(&product, id).Error

	return &product, err
}

func (r *catalogRepository) UpdateProduct(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *catalogRepository) DeleteProduct(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *catalogRepository) ReplaceProductImages(productID uint, images []models.ProductImage) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id = ?", productID).Delete(&models.ProductImage{}).Error; err != nil {
			return err
		}

		if len(images) == 0 {
			return nil
		}

		for index := range images {
			images[index].ProductID = productID
		}

		return tx.Create(&images).Error
	})
}

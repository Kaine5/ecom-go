package repository

import (
	"context"
	"errors"

	"ecom-go/internal/models"

	"gorm.io/gorm"
)

// ProductRepository defines the interface for product data access
type ProductRepo struct {
	db *gorm.DB
}

// NewProductRepo creates a new product repository
func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{
		db: db,
	}
}

// Create adds a new product to the database
func (r *ProductRepo) Create(ctx context.Context, product *models.Product) error {
	result := r.db.WithContext(ctx).Create(product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetByID retrieves a product by ID
func (r *ProductRepo) GetByID(ctx context.Context, id int) (*models.Product, error) {
	var product models.Product
	result := r.db.WithContext(ctx).First(&product, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &product, nil
}

// List retrieves all products
func (r *ProductRepo) List(ctx context.Context) ([]*models.Product, error) {
	var products []*models.Product
	result := r.db.WithContext(ctx).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

// Update updates an existing product
func (r *ProductRepo) Update(ctx context.Context, product *models.Product) error {
	result := r.db.WithContext(ctx).Save(product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete removes a product from the database
func (r *ProductRepo) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&models.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Count returns the total number of products
func (r *ProductRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&models.Product{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

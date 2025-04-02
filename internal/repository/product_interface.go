package repository

import (
	"context"

	"ecom-go/internal/models"
)

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	// Create adds a new product to the database
	Create(ctx context.Context, product *models.Product) error

	// GetByID retrieves a product by ID
	GetByID(ctx context.Context, id int) (*models.Product, error)

	// List retrieves all products
	List(ctx context.Context) ([]*models.Product, error)

	// Update updates an existing product
	Update(ctx context.Context, product *models.Product) error

	// Delete removes a product from the database
	Delete(ctx context.Context, id int) error

	// Count returns the total number of products
	Count(ctx context.Context) (int64, error)

	// ListWithPagination retrieves products with pagination
	ListWithPagination(ctx context.Context, offset, limit int) ([]*models.Product, error)
}

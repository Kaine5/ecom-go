package repository

import (
	"context"

	"ecom-go/internal/models"
)

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	// Create adds a new order to the database
	Create(ctx context.Context, order *models.Order) error

	// GetByID retrieves an order by ID
	GetByID(ctx context.Context, id int) (*models.Order, error)

	// List retrieves orders with pagination
	List(ctx context.Context, offset, limit int) ([]*models.Order, error)
}

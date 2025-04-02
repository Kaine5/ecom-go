package repository

import (
	"context"
	"errors"

	"ecom-go/internal/models"

	"gorm.io/gorm"
)

type OrderRepo struct {
	db *gorm.DB
}

// NewOrderRepo creates a new order repository
func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{
		db: db,
	}
}

// Create adds a new order to the database
func (r *OrderRepo) Create(ctx context.Context, order *models.Order) error {
	result := r.db.WithContext(ctx).Create(order)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetByID retrieves an order by ID
func (r *OrderRepo) GetByID(ctx context.Context, id int) (*models.Order, error) {
	var order models.Order
	result := r.db.WithContext(ctx).First(&order, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &order, nil
}

// List retrieves orders with pagination
func (r *OrderRepo) List(ctx context.Context, offset, limit int) ([]*models.Order, error) {
	var orders []*models.Order
	result := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}

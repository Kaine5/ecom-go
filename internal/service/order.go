package service

import (
	"context"
	"ecom-go/internal/dtos"
	"ecom-go/internal/models"
	"ecom-go/internal/repository"
	appError "ecom-go/pkg/errors"
)

type OrderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(ctx context.Context, createOrderDTO *dtos.CreateOrderDTO) (*models.Order, error) {
	order := &models.Order{
		Products: createOrderDTO.Products,
		UserID:   createOrderDTO.UserID,
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, appError.NewServerError("Failed to create order", err)
	}

	return order, nil
}

// ListOrders retrieves a list of orders with pagination
func (s *OrderService) ListOrders(ctx context.Context, offset, limit int) ([]*models.Order, error) {
	orders, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, appError.NewServerError("Failed to list orders", err)
	}

	return orders, nil
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(ctx context.Context, id int) (*models.Order, error) {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appError.NewServerError("Failed to get order", err)
	}

	return order, nil
}

package dtos

import (
	"ecom-go/internal/models"
)

// CreateOrderDTO represents the input for creating a new order
type CreateOrderDTO struct {
	UserID   int                `json:"user_id" binding:"required"`
	Products []models.OrderItem `json:"products" binding:"required"`
}

type ViewOrderDTO struct {
	ID int `json:"id"`
}

type ListOrderDTO struct {
	UserID int    `json:"user_id"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Filter string `json:"filter"`
	Sort   string `json:"sort"`
}

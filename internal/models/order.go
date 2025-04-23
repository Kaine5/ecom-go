package models

import "time"

type OrderItem struct {
	OrderItemID    int `json:"order_item_id" gorm:"uniqueIndex;primaryKey;autoIncrement"`
	ProductID      int `json:"product_id"`
	RelatedOrderID int `json:"order_id"`
	Quantity       int `json:"quantity"`
}

type Order struct {
	OrderID    int         `json:"id" gorm:"uniqueIndex;primaryKey;autoIncrement"`
	UserID     int         `json:"user_id"`
	Products   []OrderItem `json:"products" gorm:"foreignKey:RelatedOrderID"` // List of products with quantity and price
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status" gorm:"default:pending"` // e.g., "pending", "completed", "canceled"
	CreatedAt  time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

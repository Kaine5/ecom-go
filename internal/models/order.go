package models

type OrderItem struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"` // Price at the time of order
}

type Order struct {
	ID         int         `json:"id"`
	UserID     int         `json:"user_id"`
	Products   []OrderItem `json:"products"` // List of products with quantity and price
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status"` // e.g., "pending", "completed", "canceled"
	CreatedAt  string      `json:"created_at"`
	UpdatedAt  string      `json:"updated_at"`
}

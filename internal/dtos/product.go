package dtos

type ViewProductDTO struct {
	ID int `json:"id"`
}

type CreateProductDTO struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Stock       int     `json:"stock" binding:"required"`
}

type UpdateProductDTO struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

type DeleteProductDTO struct {
	ID int `json:"id"`
}

package service

import (
	"context"
	"ecom-go/internal/dtos"
	"ecom-go/internal/models"
	"ecom-go/internal/repository"
	appError "ecom-go/pkg/errors"
)

// ProductService handles business logic related to products
type ProductService struct {
	repo repository.ProductRepository
}

// NewProductService creates a new product service
func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, createProductDTO dtos.CreateProductDTO) (*models.Product, error) {
	product := &models.Product{
		Name:        createProductDTO.Name,
		Description: createProductDTO.Description,
		Price:       createProductDTO.Price,
		Stock:       createProductDTO.Stock,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, appError.NewServerError("error creating product", err)
	}

	return product, nil
}

// ListProducts retrieves all products
func (s *ProductService) ListProducts(ctx context.Context) ([]*models.Product, error) {
	products, err := s.repo.List(ctx)
	if err != nil {
		return nil, appError.NewServerError("error listing products", err)
	}

	return products, nil
}

// UpdateProduct updates a product
func (s *ProductService) UpdateProduct(ctx context.Context, id int, updateProductDTO dtos.UpdateProductDTO) (*models.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appError.NewNotFoundError("product not found")
	}

	product.Name = updateProductDTO.Name
	product.Description = updateProductDTO.Description
	product.Price = updateProductDTO.Price
	product.Stock = updateProductDTO.Stock

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, appError.NewServerError("error updating product", err)
	}

	return product, nil
}

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return appError.NewServerError("error deleting product", err)
	}

	return nil
}

// ViewProduct retrieves a product by ID
func (s *ProductService) ViewProduct(ctx context.Context, id int) (*models.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appError.NewNotFoundError("product not found")
	}

	return product, nil
}

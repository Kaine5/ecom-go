package repository

import (
	"context"

	"ecom-go/internal/models"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create adds a new user to the database
	Create(ctx context.Context, user *models.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uint) (*models.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*models.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *models.User) error

	// Delete removes a user from the database
	Delete(ctx context.Context, id uint) error

	// List retrieves users with pagination
	List(ctx context.Context, offset, limit int) ([]*models.User, error)

	// Count returns the total number of users
	Count(ctx context.Context) (int64, error)
}

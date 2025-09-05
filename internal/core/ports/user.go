package ports

import (
	"context"
	"go-ecommerce/internal/core/domain"

	"github.com/google/uuid"
)

// Methods that interacting with database
type UserRepository interface {
	// CreateUser inserts a new user into the database
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	// GetUserByID selects a user by id
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	// GetUserByEmail selects a user by email
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	// ListUsers selects a list of users with pagination
	ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error)
	// UpdateUser updates a user
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

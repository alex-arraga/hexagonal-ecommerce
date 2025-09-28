package ports

import (
	"context"
	"go-ecommerce/internal/core/domain"

	"github.com/google/uuid"
)

// UserReposity is an interface that contains methods for interacting with the repository, which will impact the database
type UserRepository interface {
	SaveUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	ListUsers(ctx context.Context, skip, limit uint64) ([]*domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}


// UserService is an interface for interacting with user-related business logic
type UserService interface {
	SaveUser(ctx context.Context, inputs domain.SaveUserInputs) (*domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	ListUsers(ctx context.Context, skip, limit uint64) ([]*domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

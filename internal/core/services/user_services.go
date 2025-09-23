package services

import (
	"context"
	"go-ecommerce/internal/adapters/shared"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"
	"log/slog"
)

type UserService struct {
	repo   ports.UserRepository
	cache  ports.UserCacheRepository
	hasher domain.PasswordHasher
}

func NewUserService(repo ports.UserRepository, cache ports.UserCacheRepository, hasher domain.PasswordHasher) ports.UserService {
	return &UserService{
		repo:   repo,
		cache:  cache,
		hasher: hasher,
	}
}

// Register creates a new user
func (us *UserService) Register(ctx context.Context, name, email, password string, role domain.UserRole) (*domain.User, error) {
	// create user domain entity applying business rules
	u, err := domain.NewUser(name, email, password, role, us.hasher)
	if err != nil {
		return nil, err
	}

	createdUser, err := us.repo.CreateUser(ctx, u)
	if err != nil {
		if err == shared.ErrConflictingData {
			return nil, err
		}
		return nil, shared.ErrInternal
	}

	// Cache the newly created user
	err = us.cache.SetUser(ctx, createdUser)
	if err != nil {
		slog.Warn("failed tu cache newly user", "error", err)
	}

	return createdUser, nil
}

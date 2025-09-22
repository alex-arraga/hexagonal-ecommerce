package services

import (
	"context"
	"go-ecommerce/internal/adapters/shared"
	"go-ecommerce/internal/adapters/shared/encoding"
	"go-ecommerce/internal/adapters/storage/cache/redis"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"
)

type UserService struct {
	repo   ports.UserRepository
	cache  ports.CacheRepository
	hasher domain.PasswordHasher
}

func NewUserService(repo ports.UserRepository, cache ports.CacheRepository, hasher domain.PasswordHasher) ports.UserService {
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

	// generate cache key and setting with the user data
	cacheKey := redis.GenerateCacheKey("user", createdUser.ID)
	userSerialized, err := encoding.Serialize(createdUser)
	if err != nil {
		return nil, shared.ErrInternal
	}

	err = us.cache.Set(ctx, cacheKey, userSerialized, 0)
	if err != nil {
		return nil, shared.ErrInternal
	}

	// delete all users of the list and update again when impact database
	err = us.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return nil, shared.ErrInternal
	}

	return createdUser, nil
}

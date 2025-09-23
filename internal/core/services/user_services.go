package services

import (
	"context"
	"encoding/json"
	"go-ecommerce/internal/adapters/shared"
	cachekeys "go-ecommerce/internal/adapters/storage/cache/cache_keys"
	cachettl "go-ecommerce/internal/adapters/storage/cache/cache_ttl"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"
	"log/slog"
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

	// create the user in the repository
	createdUser, err := us.repo.CreateUser(ctx, u)
	if err != nil {
		if err == shared.ErrConflictingData {
			return nil, err
		}
		return nil, shared.ErrInternal
	}

	// cache the newly created user
	cacheKey := cachekeys.User(createdUser.ID.String())
	userSerialized, _ := json.Marshal(createdUser)
	err = us.cache.Set(ctx, cacheKey, userSerialized, cachettl.User)
	if err != nil {
		slog.Warn("error caching user", "user_id", createdUser.ID, "error", err)
	}

	// invalid the cached list of all users
	err = us.cache.Delete(ctx, cachekeys.AllUsers())
	if err != nil {
		slog.Warn("error invalidating list of all users", "error", err)
	}

	return createdUser, nil
}

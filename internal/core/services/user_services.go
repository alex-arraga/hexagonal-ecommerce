package services

import (
	"context"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"
	"go-ecommerce/internal/core/utils"
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
func (us *UserService) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	// create user domain entity applying business rules
	u, err := domain.NewUser(user.Email, user.Name, user.Password, user.Role, us.hasher)
	if err != nil {
		return nil, err
	}

	user, err = us.repo.CreateUser(ctx, u)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	// generate cache key and setting with the user data
	cacheKey := utils.GenerateCacheKey("user", user.ID)
	userSerialized, err := utils.Serialize(user)
	if err != nil {
		return nil, domain.ErrInternal
	}

	err = us.cache.Set(ctx, cacheKey, userSerialized, 0)
	if err != nil {
		return nil, domain.ErrInternal
	}

	// delete all users of the list and update again when impact database
	err = us.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return nil, domain.ErrInternal
	}

	return user, nil
}

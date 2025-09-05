package services

import (
	"context"
	"go-ecommerce/internal/adapters/storage/cache/redis"
	"go-ecommerce/internal/adapters/storage/database/postgres/repository"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"
	"go-ecommerce/internal/core/utils"
)

type UserService struct {
	repo  repository.UserRepo
	cache redis.Redis
}

func NewUserService(repo repository.UserRepo, cache redis.Redis) ports.UserService {
	return &UserService{
		repo:  repo,
		cache: cache,
	}
}

// Register creates a new user
func (us *UserService) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	// hashing password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, domain.ErrInternal
	}

	user.Password = hashedPassword

	user, err = us.repo.CreateUser(ctx, user)
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

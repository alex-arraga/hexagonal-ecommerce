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

	"github.com/google/uuid"
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
func (us *UserService) SaveUser(ctx context.Context, inputs domain.SaveUserInputs) (*domain.User, error) {
	var user *domain.User

	if inputs.ID == uuid.Nil {
		// create a user entity applying business rules
		inputs := domain.SaveUserInputs{
			Name:     inputs.Name,
			Email:    inputs.Email,
			Password: inputs.Password,
			Role:     inputs.Role,
		}

		// find user before update
		existingUser, err := us.repo.GetUserByEmail(ctx, *inputs.Email)
		if err != nil {
			return nil, err
		}

		// checks if exist an account with the email sent
		if existingUser != nil {
			return nil, domain.ErrEmailExist
		}

		newUser, err := domain.NewUser(inputs, us.hasher)
		if err != nil {
			return nil, err
		}
		user = newUser
	} else {
		// find user before update
		existingUser, err := us.repo.GetUserByID(ctx, inputs.ID)
		if err != nil {
			return nil, err
		}

		// update user entity applying business rules
		inputs := domain.SaveUserInputs{
			ID:       inputs.ID,
			Name:     inputs.Name,
			Email:    inputs.Email,
			Password: inputs.Password,
			Role:     inputs.Role,
		}
		existingUser.UpdateUser(inputs, us.hasher)
		user = existingUser
	}

	// create the user in the repository
	result, err := us.repo.SaveUser(ctx, user)
	if err != nil {
		if err == shared.ErrConflictingData {
			return nil, err
		}
		return nil, shared.ErrInternal
	}

	// keys of cache
	emailCacheKey := cachekeys.UserByEmail(result.Email)
	cacheKey := cachekeys.User(result.ID.String())

	// cache the saved user
	userSerialized, _ := json.Marshal(result)
	err = us.cache.Set(ctx, cacheKey, userSerialized, cachettl.User)
	if err != nil {
		slog.Warn("error caching user", "user_id", result.ID, "error", err)
	}

	err = us.cache.Set(ctx, emailCacheKey, []byte(user.ID.String()), cachettl.User)
	if err != nil {
		slog.Warn("error setting user map between email and ID in cache", "user_email", user.Email, "user_id", user.ID, "error", err)
	}

	// invalid the cached list of all users
	err = us.cache.Delete(ctx, cachekeys.AllUsers())
	if err != nil {
		slog.Warn("error invalidating list of all users", "error", err)
	}

	return result, nil
}

// GetUserByEmail implements ports.UserService.
func (us *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	emailCacheKey := cachekeys.UserByEmail(email)

	idBytes, err := us.cache.Get(ctx, emailCacheKey)
	if err == nil && len(idBytes) > 0 {
		idStr := string(idBytes)

		data, err := us.cache.Get(ctx, idStr)
		if err == nil && len(data) > 0 {
			var user domain.User
			if decodeErr := json.Unmarshal(data, &user); decodeErr != nil {
				return &user, nil
			}
		}
	}

	// if the user is not cached, find in repo and cache
	user, err := us.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	serialized, err := json.Marshal(user)
	if err != nil {
		slog.Warn("error marshaling user for cache", "error", err)
	}

	// save in cache, user and a mapping between email and user.ID
	cacheKey := cachekeys.User(user.ID.String())
	err = us.cache.Set(ctx, cacheKey, serialized, cachettl.User)
	if err != nil {
		slog.Warn("error setting user in cache", "user_id", user.ID, "error", err)
	}
	err = us.cache.Set(ctx, emailCacheKey, []byte(user.ID.String()), cachettl.User)
	if err != nil {
		slog.Warn("error setting user map between email and ID in cache", "user_email", email, "user_id", user.ID, "error", err)
	}

	return user, nil
}

// GetUserByID implements ports.UserService.
func (us *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	cacheKey := cachekeys.User(id.String())

	// check if the product exist in cache, if exist return it
	data, err := us.cache.Get(ctx, cacheKey)
	if err == nil && len(data) > 0 {
		var user domain.User
		if decodeErr := json.Unmarshal(data, &user); decodeErr != nil {
			return &user, nil
		}
	}

	// else find user in repository
	user, err := us.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// set cache
	serialized, err := json.Marshal(user)
	if err != nil {
		slog.Warn("error marshaling user for cache", "error", err)
	}

	err = us.cache.Set(ctx, cacheKey, serialized, cachettl.User)
	if err != nil {
		slog.Warn("error setting user in cache", "user_id", user.ID, "error", err)
	}

	return user, nil
}

// ListUsers implements ports.UserService.
func (us *UserService) ListUsers(ctx context.Context, skip uint64, limit uint64) ([]*domain.User, error) {
	// check if the products exists in cache
	data, err := us.cache.Get(ctx, cachekeys.AllUsers())
	if err == nil && len(data) > 0 {
		var users []*domain.User
		if decodeErr := json.Unmarshal(data, &users); decodeErr != nil {
			return users, nil
		}
	}

	// find users in repository
	users, err := us.repo.ListUsers(ctx, 0, 20)
	if err != nil {
		return nil, err
	}

	// serialize products
	usersSerialized, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}

	// regenerate list of products
	err = us.cache.Set(ctx, cachekeys.AllUsers(), usersSerialized, cachettl.User)
	if err != nil {
		slog.Warn("error caching users", "error", err)
	}

	return users, nil
}

// DeleteUser implements ports.UserService.
func (us *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	cacheKey := cachekeys.User(id.String())

	err := us.repo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	err = us.cache.Delete(ctx, cacheKey)
	if err != nil {
		slog.Warn("error deleteing user of cache", "user_id", id, "error", err)
	}

	err = us.cache.Delete(ctx, cachekeys.AllUsers())
	if err != nil {
		slog.Warn("error invalidating list of all users", "error", err)
	}

	return nil
}

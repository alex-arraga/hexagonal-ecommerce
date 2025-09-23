package redis

import (
	"context"
	"go-ecommerce/internal/adapters/shared/encoding"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewUserCache(client *redis.Client, ttl time.Duration) ports.UserCacheRepository {
	return &UserCache{
		client: client,
		ttl:    ttl,
	}
}

const userKey = "user"
const userSetKey = "user_keys_set"

// SetUser implements ports.UserCacheRepositroy.
func (uc *UserCache) SetUser(ctx context.Context, user *domain.User) error {
	key := GenerateCacheKey(userKey, user.ID)

	data, err := encoding.Serialize(user)
	if err != nil {
		return err
	}

	// Store the user key in a Redis set for easy retrieval and deletion later
	if err = uc.client.SAdd(ctx, userSetKey, key).Err(); err != nil {
		return err
	}

	// Set the user data with the specified TTL
	return uc.client.Set(ctx, key, data, uc.ttl).Err()
}

// GetUser implements ports.UserCacheRepositroy.
func (uc *UserCache) GetUser(ctx context.Context, id string) (*domain.User, error) {
	key := GenerateCacheKey(userKey, id)
	data, err := uc.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var user domain.User
	if err := encoding.Deserialize(data, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteUser implements ports.UserCacheRepositroy.
func (uc *UserCache) DeleteUser(ctx context.Context, id string) error {
	key := GenerateCacheKey(userKey, id)

	if err := uc.client.SRem(ctx, userSetKey, key).Err(); err != nil {
		if err == redis.Nil {
			return nil // Key does not exist in the set
		}
		return err
	}

	return uc.client.Del(ctx, key).Err()
}

// DeleteAllUsers implements ports.UserCacheRepositroy.
func (uc *UserCache) DeleteAllUsers(ctx context.Context) error {
	keys, err := uc.client.SMembers(ctx, userSetKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil // Set does not exist
		}
	}

	if len(keys) > 0 {
		if err := uc.client.Del(ctx, keys...).Err(); err != nil {
			return err
		}
	}

	// Clear the set of user keys
	return uc.client.Del(ctx, userSetKey).Err()
}

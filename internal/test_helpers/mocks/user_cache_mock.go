package mocks

import (
	"context"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"
)

type UserCacheMock struct {
	store map[string]string
}

func NewUserCacheMock() ports.UserCacheRepository {
	return &UserCacheMock{store: make(map[string]string)}
}

// SetUser implements ports.UserCacheRepository.
func (um *UserCacheMock) SetUser(ctx context.Context, user *domain.User) error {
	return nil
}

// GetUser implements ports.UserCacheRepository.
func (um *UserCacheMock) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return nil, nil
}

// DeleteUser implements ports.UserCacheRepository.
func (um *UserCacheMock) DeleteUser(ctx context.Context, id string) error {
	return nil
}

// DeleteAllUsers implements ports.UserCacheRepository.
func (um *UserCacheMock) DeleteAllUsers(ctx context.Context) error {
	return nil
}

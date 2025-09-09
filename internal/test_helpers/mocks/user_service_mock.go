package mocks

import (
	"context"
	"go-ecommerce/internal/core/domain"
)

type MockUserService struct {
	RegisterFunc func(ctx context.Context, user *domain.User) (*domain.User, error)
}

func (m *MockUserService) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	return m.RegisterFunc(ctx, user)
}

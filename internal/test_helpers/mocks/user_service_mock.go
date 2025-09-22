package mocks

import (
	"context"
	"go-ecommerce/internal/core/domain"
)

type MockUserService struct {
	RegisterFunc func(ctx context.Context, name, email, password string, role domain.UserRole) (*domain.User, error)
}

func (m *MockUserService) Register(ctx context.Context, name, email, password string, role domain.UserRole) (*domain.User, error) {
	return m.RegisterFunc(ctx, name, email, password, role)
}

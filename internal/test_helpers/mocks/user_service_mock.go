package mocks

import (
	"context"
	"go-ecommerce/internal/core/domain"

	"github.com/google/uuid"
)

type MockUserService struct {
	SaveFunc func(ctx context.Context, inputs domain.SaveUserInputs) (*domain.User, error)
}

// SaveUser implements ports.UserService.
func (m *MockUserService) SaveUser(ctx context.Context, inputs domain.SaveUserInputs) (*domain.User, error) {
	return m.SaveFunc(ctx, inputs)
}

// GetUserByEmail implements ports.UserService.
func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	panic("unimplemented")
}

// GetUserByID implements ports.UserService.
func (m *MockUserService) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	panic("unimplemented")
}

// ListUsers implements ports.UserService.
func (m *MockUserService) ListUsers(ctx context.Context, skip uint64, limit uint64) ([]*domain.User, error) {
	panic("unimplemented")
}

// DeleteUser implements ports.UserService.
func (m *MockUserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	panic("unimplemented")
}

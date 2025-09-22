package ports

import (
	"context"
	"go-ecommerce/internal/core/domain"
)

// interface that category_repository implements
type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	GetCategoryByID(ctx context.Context, id uint64) (*domain.Category, error)
	ListCategories(ctx context.Context) ([]*domain.Category, error)
	UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id uint64) error
}

// interface that category_service implements
type CategoryService interface {
	RegisterCategory(ctx context.Context, name string) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id uint64) error
}

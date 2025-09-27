package ports

import (
	"context"
	"go-ecommerce/internal/core/domain"
)

// interface that category_repository implements
type CategoryRepository interface {
	SaveCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	GetCategoryByID(ctx context.Context, id uint64) (*domain.Category, error)
	ListCategories(ctx context.Context) ([]*domain.Category, error)
	DeleteCategory(ctx context.Context, id uint64) error
}

// interface that category_service implements
type CategoryService interface {
	SaveCategory(ctx context.Context, id uint64, name string) (*domain.Category, error)
	GetCategoryByID(ctx context.Context, id uint64) (*domain.Category, error)
	ListCategories(ctx context.Context) ([]*domain.Category, error)
	DeleteCategory(ctx context.Context, id uint64) error
}

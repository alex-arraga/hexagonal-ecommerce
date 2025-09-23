package ports

import (
	"context"
	"go-ecommerce/internal/core/domain"

	"github.com/google/uuid"
)

// ProductRepository is an interface that contains methods for interacting with the repository, which will impact the database
type ProductRepository interface {
	SaveProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetProductById(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	ListProducts(ctx context.Context) ([]*domain.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

// ProductService is an interface for interacting with product-related business logic
type SaveProductInputs struct {
	ID         uuid.UUID
	Name       string
	Image      string
	SKU        string
	Price      float64
	Stock      int64
	CategoryID uint64
}
type ProductService interface {
	SaveProduct(ctx context.Context, inputs SaveProductInputs) (*domain.Product, error)
	GetProductById(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	ListProducts(ctx context.Context) ([]*domain.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

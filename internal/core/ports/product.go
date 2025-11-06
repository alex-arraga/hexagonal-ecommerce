package ports

import (
	"context"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports/ports_dtos"

	"github.com/google/uuid"
)

// ProductRepository is an interface that contains methods for interacting with the repository, which will impact the database
type ProductRepository interface {
	SaveProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetProductById(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	ListProducts(ctx context.Context) ([]*domain.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

type ProductService interface {
	SaveProduct(ctx context.Context, inputs ports_dtos.SaveProductInputs) (*domain.Product, error)
	GetProductById(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	ListProducts(ctx context.Context) ([]*domain.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

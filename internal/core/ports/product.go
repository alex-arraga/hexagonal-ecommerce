package ports

import (
	"context"
	"go-ecommerce/internal/core/domain"

	"github.com/google/uuid"
)

type ProductRepository interface {
	SaveProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetProductById(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	ListProducts(ctx context.Context) ([]*domain.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

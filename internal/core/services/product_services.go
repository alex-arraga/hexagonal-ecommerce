package services

import (
	"context"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"

	"github.com/google/uuid"
)

type ProductService struct {
	repo  ports.ProductRepository
	cache ports.CacheRepository
}

func NewProductService(repo ports.ProductRepository, cache ports.CacheRepository) ports.ProductService {
	return &ProductService{
		repo:  repo,
		cache: cache,
	}
}

// SaveProduct implements ports.ProductService.
func (p *ProductService) SaveProduct(ctx context.Context, inputs ports.SaveProductInputs) (*domain.Product, error) {
	panic("unimplemented")
}

// GetProductById implements ports.ProductService.
func (p *ProductService) GetProductById(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	panic("unimplemented")
}

// ListProducts implements ports.ProductService.
func (p *ProductService) ListProducts(ctx context.Context) ([]*domain.Product, error) {
	panic("unimplemented")
}

// DeleteProduct implements ports.ProductService.
func (p *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	panic("unimplemented")
}

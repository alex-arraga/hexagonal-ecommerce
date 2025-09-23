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
func (ps *ProductService) SaveProduct(ctx context.Context, inputs ports.SaveProductInputs) (*domain.Product, error) {
	var product *domain.Product

	if inputs.ID == uuid.Nil {
		// create a new product if inputs.ID doesn't exist
		newProduct, err := domain.NewProduct(
			inputs.Name,
			inputs.SKU,
			inputs.Image,
			inputs.Stock,
			inputs.Price,
			inputs.CategoryID,
		)
		if err != nil {
			return nil, err
		}
		product = newProduct
		
	} else {
		product, err := ps.repo.GetProductById(ctx, inputs.ID)
		if err != nil {
			return nil, err
		}

		product.Update(
			inputs.Name,
			inputs.SKU,
			inputs.Image,
			inputs.Stock,
			inputs.Price,
			inputs.CategoryID,
		)
	}

	result, err := ps.repo.SaveProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetProductById implements ports.ProductService.
func (ps *ProductService) GetProductById(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	panic("unimplemented")
}

// ListProducts implements ports.ProductService.
func (ps *ProductService) ListProducts(ctx context.Context) ([]*domain.Product, error) {
	panic("unimplemented")
}

// DeleteProduct implements ports.ProductService.
func (ps *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	panic("unimplemented")
}

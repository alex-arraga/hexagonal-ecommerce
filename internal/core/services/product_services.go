package services

import (
	"context"
	"go-ecommerce/internal/adapters/shared/encoding"
	cachekeys "go-ecommerce/internal/adapters/storage/cache/cache_keys"
	cachettl "go-ecommerce/internal/adapters/storage/cache/cache_ttl"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"
	"log/slog"

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

	createdProduct, err := ps.repo.SaveProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	// create new product cache key and serialize product created or udpated
	cacheKey := cachekeys.Product(createdProduct.ID.String())
	productSerialized, err := encoding.Serialize(createdProduct)
	if err != nil {
		return nil, err
	}

	// set product in cache
	err = ps.cache.Set(ctx, cacheKey, productSerialized, cachettl.Product)
	if err != nil {
		slog.Warn("error caching new created product", "product_id", createdProduct.ID, "error", err)
	}

	// invalidate product list
	err = ps.cache.Delete(ctx, cachekeys.AllProducts())
	if err != nil {
		slog.Warn("error invalidating list of all products", "error", err)
	}

	return createdProduct, nil
}

// GetProductById implements ports.ProductService.
func (ps *ProductService) GetProductById(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	cacheKey := cachekeys.Product(id.String())

	// check if the product exist in cache, if exist return it
	val, err := ps.cache.Get(ctx, cacheKey)
	if err == nil {
		var product domain.Product
		if decodeErr := encoding.Deserialize(val, &product); decodeErr != nil {
			return &product, nil
		}
	}

	// else find product in repository
	p, err := ps.repo.GetProductById(ctx, id)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// ListProducts implements ports.ProductService.
func (ps *ProductService) ListProducts(ctx context.Context) ([]*domain.Product, error) {
	panic("unimplemented")
}

// DeleteProduct implements ports.ProductService.
func (ps *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	panic("unimplemented")
}

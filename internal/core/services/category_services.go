package services

import (
	"context"
	"go-ecommerce/internal/adapters/shared"
	"go-ecommerce/internal/adapters/shared/encoding"
	cachekeys "go-ecommerce/internal/adapters/storage/cache/cache_keys"
	cachettl "go-ecommerce/internal/adapters/storage/cache/cache_ttl"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"
	"log/slog"
)

type CategoryService struct {
	repo  ports.CategoryRepository
	cache ports.CacheRepository
}

func NewCategoryService(repo ports.CategoryRepository, cache ports.CacheRepository) ports.CategoryService {
	return &CategoryService{
		repo:  repo,
		cache: cache,
	}
}

// RegisterCategory implements ports.CategoryService.
func (cs *CategoryService) RegisterCategory(ctx context.Context, name string) (*domain.Category, error) {
	domainCategory, err := domain.NewCategory(name)
	if err != nil {
		return nil, err
	}

	createdCategory, err := cs.repo.CreateCategory(ctx, domainCategory)
	if err != nil {
		return nil, err
	}

	// generate cache key and setting with the user data
	cacheKey := cachekeys.Category(createdCategory.ID)
	categorySerialized, err := encoding.Serialize(createdCategory)
	if err != nil {
		return nil, shared.ErrInternal
	}

	err = cs.cache.Set(ctx, cacheKey, categorySerialized, cachettl.Category)
	if err != nil {
		slog.Warn("error caching category", "category_id", createdCategory.ID, "error", err)
	}

	// invalid the cached list of all categories
	err = cs.cache.Delete(ctx, cachekeys.AllCategories())
	if err != nil {
		slog.Warn("error invalidating list of all categories", "error", err)
	}

	return createdCategory, nil
}

// DeleteCategory implements ports.CategoryService.
func (cs *CategoryService) DeleteCategory(ctx context.Context, id uint64) error {
	err := cs.repo.DeleteCategory(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

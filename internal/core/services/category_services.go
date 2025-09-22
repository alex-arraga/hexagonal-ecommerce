package services

import (
	"context"
	"go-ecommerce/internal/adapters/shared"
	"go-ecommerce/internal/adapters/shared/encoding"
	"go-ecommerce/internal/adapters/storage/cache/redis"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"
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
	cacheKey := redis.GenerateCacheKey("categories", createdCategory.ID)
	categorySerialized, err := encoding.Serialize(createdCategory)
	if err != nil {
		return nil, shared.ErrInternal
	}

	err = cs.cache.Set(ctx, cacheKey, categorySerialized, 10)
	if err != nil {
		return nil, shared.ErrInternal
	}

	// delete all users of the list and update again when impact database
	err = cs.cache.DeleteByPrefix(ctx, "categories:*")
	if err != nil {
		return nil, shared.ErrInternal
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

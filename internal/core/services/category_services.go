package services

import (
	"context"
	"encoding/json"
	"go-ecommerce/internal/adapters/shared"
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
func (cs *CategoryService) SaveCategory(ctx context.Context, id uint64, name string) (*domain.Category, error) {
	var category *domain.Category

	if id == 0 {
		newCategory, err := domain.NewCategory(name)
		if err != nil {
			return nil, err
		}
		category = newCategory
	} else {
		category.UpdateCategory(name)
	}

	result, err := cs.repo.SaveCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	// generate cache key and setting with the user data
	cacheKey := cachekeys.Category(result.ID)
	categorySerialized, err := json.Marshal(result)
	if err != nil {
		return nil, shared.ErrInternal
	}

	// set the category in cache
	_ = cs.cache.Set(ctx, cacheKey, categorySerialized, cachettl.Category)

	// invalid the cached list of all categories
	_ = cs.cache.Delete(ctx, cachekeys.AllCategories())

	return result, nil
}

// GetCategoryByID implements ports.CategoryService.
func (cs *CategoryService) GetCategoryByID(ctx context.Context, id uint64) (*domain.Category, error) {
	cacheKey := cachekeys.Category(id)

	// validate if category is saved in cache
	data, err := cs.cache.Get(ctx, cacheKey)
	if err == nil && len(data) > 0 {
		var category *domain.Category
		if decodeErr := json.Unmarshal(data, &category); decodeErr == nil {
			return category, nil
		}
	}

	// else find category in repository and set in cache
	category, err := cs.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// set cache
	serialized, err := json.Marshal(category)
	if err != nil {
		slog.Warn("Error marshaling category for cache", "error", err)
	}

	_ = cs.cache.Set(ctx, cacheKey, serialized, cachettl.Category)

	return category, nil
}

// ListCategories implements ports.CategoryService.
func (cs *CategoryService) ListCategories(ctx context.Context) ([]*domain.Category, error) {
	// validate if category is saved in cache
	data, err := cs.cache.Get(ctx, cachekeys.AllCategories())
	if err == nil && len(data) > 0 {
		var categories []*domain.Category
		if decodeErr := json.Unmarshal(data, &categories); decodeErr == nil {
			return categories, nil
		}
	}

	// else find categories in repository
	categories, err := cs.repo.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	// save in cache
	serialized, err := json.Marshal(categories)
	if err != nil {
		slog.Warn("Error marshaling categories for cache", "error", err)
	}

	// regenerate list of categories
	_ = cs.cache.Set(ctx, cachekeys.AllCategories(), serialized, cachettl.Category)

	return categories, nil
}

// DeleteCategory implements ports.CategoryService.
func (cs *CategoryService) DeleteCategory(ctx context.Context, id uint64) error {
	err := cs.repo.DeleteCategory(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

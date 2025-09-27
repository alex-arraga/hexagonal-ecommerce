package repository

import (
	"context"
	"go-ecommerce/internal/adapters/storage/database/postgres/database_dtos"
	"go-ecommerce/internal/adapters/storage/database/postgres/models"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"

	"gorm.io/gorm"
)

type CategoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) ports.CategoryRepository {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) SaveCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	categoryDb := database_dtos.ConvertCategoryDomainToModel(category)

	if category.ID != 0 {
		if result := r.db.WithContext(ctx).Where("id = ?", category.ID).Updates(categoryDb); result.Error != nil {
			if result.RowsAffected == 0 {
				return nil, domain.ErrCategoryNotFound
			}
			return nil, result.Error
		}
	} else {
		if result := r.db.WithContext(ctx).Create(categoryDb); result.Error != nil {
			return nil, result.Error
		}
	}

	domainCategory := database_dtos.ConvertCategoryModelToDomain(categoryDb)
	return domainCategory, nil
}

func (r *CategoryRepo) GetCategoryByID(ctx context.Context, id uint64) (*domain.Category, error) {
	var categoryDb *models.CategoryModel

	if result := r.db.WithContext(ctx).Where("id = ?").First(categoryDb); result.Error != nil {
		if result.RowsAffected == 0 {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, result.Error
	}

	domainCategory := database_dtos.ConvertCategoryModelToDomain(categoryDb)
	return domainCategory, nil
}

func (r *CategoryRepo) ListCategories(ctx context.Context) ([]*domain.Category, error) {
	var categoriesDb []*models.CategoryModel

	if result := r.db.WithContext(ctx).Find(categoriesDb); result.Error != nil {
		return nil, result.Error
	}

	domainCategories := database_dtos.ConvertCategoriesModelsToDomain(categoriesDb)
	return domainCategories, nil
}

func (r *CategoryRepo) DeleteCategory(ctx context.Context, id uint64) error {
	var categoryDb *models.CategoryModel

	if result := r.db.WithContext(ctx).Where("id = ?").Delete(categoryDb); result.Error != nil {
		return result.Error
	}
	return nil
}

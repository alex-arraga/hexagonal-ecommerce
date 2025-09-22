package database_dtos

import (
	"go-ecommerce/internal/adapters/storage/database/postgres/models"
	"go-ecommerce/internal/core/domain"
)

// domain.User -> DB model
func ConvertCategoryDomainToModel(category *domain.Category) *models.CategoryModel {
	return &models.CategoryModel{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

// domain.Users -> DB models
func ConvertCategoriesDomainToModels(categories []*domain.Category) []*models.CategoryModel {
	var categoriesModels []*models.CategoryModel

	for _, c := range categories {
		categoriesModels = append(categoriesModels, &models.CategoryModel{
			ID:        c.ID,
			Name:      c.Name,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		})
	}

	return categoriesModels
}

// DB model -> domain.User
func ConvertCategoryModelToDomain(categoryModel *models.CategoryModel) *domain.Category {
	return &domain.Category{
		ID:        categoryModel.ID,
		Name:      categoryModel.Name,
		CreatedAt: categoryModel.CreatedAt,
		UpdatedAt: categoryModel.UpdatedAt,
	}
}

// DB models -> domain.Users
func ConvertCategoriesModelsToDomain(categoryModels []*models.CategoryModel) []*domain.Category {
	var domainCategories []*domain.Category

	for _, c := range categoryModels {
		domainCategories = append(domainCategories, &domain.Category{
			ID:        c.ID,
			Name:      c.Name,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		})
	}

	return domainCategories
}

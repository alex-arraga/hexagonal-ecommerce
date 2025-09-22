package database_dtos

import (
	"go-ecommerce/internal/adapters/storage/database/postgres/models"
	"go-ecommerce/internal/core/domain"
)

// domain.User -> DB model
func ConvertProductDomainToModel(p *domain.Product) *models.ProductModel {
	return &models.ProductModel{
		ID:         p.ID,
		Name:       p.Name,
		SKU:        p.SKU,
		Stock:      p.Stock,
		Price:      p.Price,
		Image:      p.Image,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
		CategoryID: p.CategoryID,
		Category:   *ConvertCategoryDomainToModel(p.Category),
	}
}

// domain.Users -> DB models
func ConvertProductsDomainToModels(products []*domain.Product) []*models.ProductModel {
	var productsModels []*models.ProductModel

	for _, p := range products {
		productsModels = append(productsModels, &models.ProductModel{
			ID:         p.ID,
			Name:       p.Name,
			SKU:        p.SKU,
			Stock:      p.Stock,
			Price:      p.Price,
			Image:      p.Image,
			CreatedAt:  p.CreatedAt,
			UpdatedAt:  p.UpdatedAt,
			CategoryID: p.CategoryID,
			Category:   *ConvertCategoryDomainToModel(p.Category),
		})
	}

	return productsModels
}

// DB model -> domain.User
func ConvertProductModelToDomain(p *models.ProductModel) *domain.Product {
	return &domain.Product{
		ID:         p.ID,
		Name:       p.Name,
		SKU:        p.SKU,
		Stock:      p.Stock,
		Price:      p.Price,
		Image:      p.Image,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
		CategoryID: p.CategoryID,
		Category:   ConvertCategoryModelToDomain(&p.Category),
	}
}

// DB models -> domain.Users
func ConvertProductsModelsToDomain(products []*models.ProductModel) []*domain.Product {
	var productsDomain []*domain.Product

	for _, p := range products {
		productsDomain = append(productsDomain, &domain.Product{
			ID:         p.ID,
			Name:       p.Name,
			SKU:        p.SKU,
			Stock:      p.Stock,
			Price:      p.Price,
			Image:      p.Image,
			CreatedAt:  p.CreatedAt,
			UpdatedAt:  p.UpdatedAt,
			CategoryID: p.CategoryID,
			Category:   ConvertCategoryModelToDomain(&p.Category),
		})
	}

	return productsDomain
}

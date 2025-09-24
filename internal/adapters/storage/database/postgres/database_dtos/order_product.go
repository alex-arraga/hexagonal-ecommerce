package database_dtos

import (
	"go-ecommerce/internal/adapters/storage/database/postgres/models"
	"go-ecommerce/internal/core/domain"
)

// domain.OrderProduct -> DB model
func ConvertOrderProductDomainToModel(op *domain.OrderProduct) *models.OrderProductModel {
	return &models.OrderProductModel{
		ID:        op.ID,
		OrderID:   op.OrderID,
		ProductID: op.ProductID,
		Quantity:  op.Quantity,
		CreatedAt: op.CreatedAt,
		UpdatedAt: op.UpdatedAt,
	}
}

// domain.OrderProducts -> DB models
func ConvertOrderProductsDomainToModels(orderProducts []*domain.OrderProduct) []*models.OrderProductModel {
	var orderProductsModels []*models.OrderProductModel

	for _, op := range orderProducts {
		orderProductsModels = append(orderProductsModels, &models.OrderProductModel{
			ID:        op.ID,
			OrderID:   op.OrderID,
			ProductID: op.ProductID,
			Quantity:  op.Quantity,
			CreatedAt: op.CreatedAt,
			UpdatedAt: op.UpdatedAt,
		})
	}

	return orderProductsModels
}

// DB model -> domain.OrderProduct
func ConvertOrderProductModelToDomain(op *models.OrderProductModel) *domain.OrderProduct {
	return &domain.OrderProduct{
		ID:        op.ID,
		OrderID:   op.OrderID,
		ProductID: op.ProductID,
		Quantity:  op.Quantity,
		CreatedAt: op.CreatedAt,
		UpdatedAt: op.UpdatedAt,
	}
}

// DB models -> domain.OrderProducts
func ConvertOrderProductModelsToDomains(orderProducts []*models.OrderProductModel) []*domain.OrderProduct {
	var orderProductsDomain []*domain.OrderProduct

	for _, op := range orderProducts {
		orderProductsDomain = append(orderProductsDomain, &domain.OrderProduct{
			ID:        op.ID,
			OrderID:   op.OrderID,
			ProductID: op.ProductID,
			Quantity:  op.Quantity,
			CreatedAt: op.CreatedAt,
			UpdatedAt: op.UpdatedAt,
		})
	}

	return orderProductsDomain
}

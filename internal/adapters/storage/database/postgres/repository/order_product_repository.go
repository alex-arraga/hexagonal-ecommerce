package repository

import (
	"context"
	"go-ecommerce/internal/adapters/storage/database/postgres/database_dtos"
	"go-ecommerce/internal/adapters/storage/database/postgres/models"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderProductRepo struct {
	db *gorm.DB
}

func NewOrderProductRepo(db *gorm.DB) ports.OrderProductRepository {
	return &OrderProductRepo{db: db}
}

// SaveOrderProduct implements ports.OrderProductRepository.
func (opr *OrderProductRepo) SaveOrderProduct(ctx context.Context, orderProduct *domain.OrderProduct) (*domain.OrderProduct, error) {
	orderProductDb := database_dtos.ConvertOrderProductDomainToModel(orderProduct)

	// if exist order.ID update, else create new order
	if orderProduct.ID != uuid.Nil {
		if result := opr.db.WithContext(ctx).Where("id = ?", orderProduct.ID).Updates(orderProductDb); result.Error != nil {
			if result.RowsAffected == 0 {
				return nil, domain.ErrOrderProductNotFound
			}
			return nil, result.Error
		}
	} else {
		if result := opr.db.WithContext(ctx).Create(orderProductDb); result.Error != nil {
			return nil, result.Error
		}
	}

	orderProductDomain := database_dtos.ConvertOrderProductModelToDomain(orderProductDb)
	return orderProductDomain, nil
}

// GetOrderProductById implements ports.OrderProductRepository.
func (opr *OrderProductRepo) GetOrderProductById(ctx context.Context, id uuid.UUID) (*domain.OrderProduct, error) {
	var orderProductDb *models.OrderProductModel

	if result := opr.db.WithContext(ctx).First(orderProductDb, "id = ?", id); result.Error != nil {
		if result.RowsAffected == 0 {
			return nil, domain.ErrOrderProductNotFound
		}
		return nil, result.Error
	}

	orderProductDomain := database_dtos.ConvertOrderProductModelToDomain(orderProductDb)
	return orderProductDomain, nil
}

// ListOrderProducts implements ports.OrderProductRepository.
func (opr *OrderProductRepo) ListOrderProducts(ctx context.Context, orderID uuid.UUID) ([]*domain.OrderProduct, error) {
	var orderProductDb []*models.OrderProductModel

	query := opr.db.WithContext(ctx)
	if orderID != uuid.Nil {
		query = query.Where("order_id = ?", orderID)
	}

	result := query.Find(&orderProductDb)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return nil, domain.ErrOrdersProductNotFound
		}
		return nil, result.Error
	}

	orderProductDomain := database_dtos.ConvertOrderProductModelsToDomains(orderProductDb)
	return orderProductDomain, nil
}

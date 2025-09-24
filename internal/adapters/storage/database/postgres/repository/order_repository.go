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

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) ports.OrderRepository {
	return &OrderRepo{db: db}
}

// SaveOrder implements ports.OrderRepository.
func (or *OrderRepo) SaveOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	orderDb := database_dtos.ConvertOrderDomainToModel(order)

	// if exist product.ID update, else create new product
	if order.ID != uuid.Nil {
		if result := or.db.WithContext(ctx).Model(orderDb).Where("id = ?", order.ID).Updates(order); result.Error != nil {
			if result.RowsAffected == 0 {
				return nil, domain.ErrProductNotFound
			}
			return nil, result.Error
		}
	} else {
		if result := or.db.WithContext(ctx).Create(orderDb); result.Error != nil {
			return nil, result.Error
		}
	}

	orderDomain := database_dtos.ConvertOrderModelToDomain(orderDb)
	return orderDomain, nil
}

// GetOrderById implements ports.OrderRepository.
func (or *OrderRepo) GetOrderById(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	var orderDb *models.OrderModel

	if result := or.db.WithContext(ctx).First(orderDb, "id = ?", id); result.Error != nil {
		if result.RowsAffected == 0 {
			return nil, domain.ErrProductNotFound
		}
		return nil, result.Error
	}

	orderDomain := database_dtos.ConvertOrderModelToDomain(orderDb)
	return orderDomain, nil
}

// ListOrders implements ports.OrderRepository.
func (or *OrderRepo) ListOrders(ctx context.Context) ([]*domain.Order, error) {
	var orderDb []*models.OrderModel

	if result := or.db.WithContext(ctx).Find(orderDb); result.Error != nil {
		return nil, result.Error
	}

	orderDomain := database_dtos.ConvertOrdersModelsToDomain(orderDb)
	return orderDomain, nil
}

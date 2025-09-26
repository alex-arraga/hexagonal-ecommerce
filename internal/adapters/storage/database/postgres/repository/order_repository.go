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
	cart ports.CartService
	opr  ports.OrderProductRepository
	db   *gorm.DB
}

func NewOrderRepo(db *gorm.DB) ports.OrderRepository {
	return &OrderRepo{db: db}
}

// SaveOrder implements ports.OrderRepository.
func (or *OrderRepo) SaveOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	cart, err := or.cart.GetCart(ctx, order.UserID)
	if err != nil {
		return nil, err
	}

	orderDb := database_dtos.ConvertOrderDomainToModel(order)

	// if exist order.ID update, else create new order
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

	// creates order-product for each item of cart
	for _, item := range cart.Items {
		op := domain.NewOrderProduct(order.ID, item.ProductID, item.Quantity)
		_, err := or.opr.SaveOrderProduct(ctx, op)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, *op)
	}

	orderDomain := database_dtos.ConvertOrderModelToDomain(orderDb)
	return orderDomain, nil
}

// GetOrderById implements ports.OrderRepository.
func (or *OrderRepo) GetOrderById(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	var orderDb *models.OrderModel

	if result := or.db.WithContext(ctx).Preload("Items").First(orderDb, "id = ?", id); result.Error != nil {
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

	if result := or.db.WithContext(ctx).Preload("Items").Find(orderDb); result.Error != nil {
		return nil, result.Error
	}

	orderDomain := database_dtos.ConvertOrdersModelsToDomain(orderDb)
	return orderDomain, nil
}

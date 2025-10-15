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
	ops  ports.OrderProductService
	db   *gorm.DB
}

func NewOrderRepo(cart ports.CartService, opr ports.OrderProductService, db *gorm.DB) ports.OrderRepository {
	return &OrderRepo{
		cart: cart,
		ops:  opr,
		db:   db,
	}
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
		if result := or.db.WithContext(ctx).Preload("User").Preload("Items").Where("id = ?", order.ID).Updates(orderDb); result.Error != nil {
			if result.RowsAffected == 0 {
				return nil, domain.ErrOrderNotFound
			}
			return nil, result.Error
		}
	} else {
		if result := or.db.WithContext(ctx).Preload("User").Preload("Items").Create(orderDb); result.Error != nil {
			return nil, result.Error
		}
	}

	// creates order-product for each item of cart
	for _, item := range cart.Items {
		op, err := or.ops.AddProductToOrder(ctx, orderDb.ID, item.ProductID, item.Quantity)
		if err != nil {
			return nil, err
		}
		opModel := database_dtos.ConvertOrderProductDomainToModel(op)
		orderDb.Items = append(orderDb.Items, *opModel)
	}

	orderDomain := database_dtos.ConvertOrderModelToDomain(orderDb)
	return orderDomain, nil
}

// GetOrderById implements ports.OrderRepository.
func (or *OrderRepo) GetOrderById(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	var orderDb = &models.OrderModel{}

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

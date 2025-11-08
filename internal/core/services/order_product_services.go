package services

import (
	"context"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"

	"github.com/google/uuid"
)

type OrderProductService struct {
	repo ports.OrderProductRepository
}

func NewOrderProductService(repo ports.OrderProductRepository) ports.OrderProductService {
	return &OrderProductService{repo: repo}
}

// AddProductToOrder implements ports.OrderProductService.
func (ops *OrderProductService) AddProductToOrder(ctx context.Context, orderID, productID uuid.UUID, quantity int16) (*domain.OrderProduct, error) {
	orderProduct := domain.NewOrderProduct(orderID, productID, quantity)
	savedOrderProduct, err := ops.repo.SaveOrderProduct(ctx, orderProduct)
	if err != nil {
		return nil, err
	}
	return savedOrderProduct, nil
}

// GetByOrderID implements ports.OrderProductService.
func (ops *OrderProductService) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*domain.OrderProduct, error) {
	items, err := ops.repo.ListOrderProducts(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

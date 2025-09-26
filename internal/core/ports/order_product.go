package ports

import (
	"context"
	"go-ecommerce/internal/core/domain"

	"github.com/google/uuid"
)

type OrderProductRepository interface {
	SaveOrderProduct(ctx context.Context, orderProduct *domain.OrderProduct) (*domain.OrderProduct, error)
	GetOrderProductById(ctx context.Context, id uuid.UUID) (*domain.OrderProduct, error)
	ListOrderProducts(ctx context.Context) ([]*domain.OrderProduct, error)
}

type OrderProductService interface {
	AddProductToOrder(ctx context.Context, orderID, productID uuid.UUID, quantity uint8) (*domain.OrderProduct, error)
}

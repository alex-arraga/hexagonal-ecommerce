package ports

import (
	"context"
	"go-ecommerce/internal/core/domain"

	"github.com/google/uuid"
)

type OrderRepository interface {
	SaveOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)
	GetOrderById(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	ListOrders(ctx context.Context) ([]*domain.Order, error)
}

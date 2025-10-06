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

// SaveOrderInputs is the input struct for saving or updating an order
type SaveOrderInputs struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	Currency          domain.Currencies
	SubTotal          float64
	Disscount         *float64
	DisscountTypes    *domain.DisscountTypes
	ExternalReference *string
	PaymentID         *string
	PayStatus         domain.PayStatus
	PayStatusDetail   *domain.PayStatusDetail
}

// OrderService is an interface for interacting with order-related business logic
type OrderService interface {
	SaveOrder(ctx context.Context, inputs SaveOrderInputs) (*domain.Order, error)
	GetOrderById(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	ListOrders(ctx context.Context) ([]*domain.Order, error)
}

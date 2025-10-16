package ports

import (
	"context"
	"go-ecommerce/internal/core/domain"

	"github.com/google/uuid"
)

type Amount struct {
	SubTotal   float64
	Discount float64
	Total      float64
}

type CartService interface {
	GetCart(ctx context.Context, userId uuid.UUID) (*domain.Cart, error)
	AddItemToCart(ctx context.Context, userId, productId uuid.UUID, quantity int16) error
	CalcItemsAmount(ctx context.Context, userId uuid.UUID) (*Amount, error)
	RemoveItem(ctx context.Context, userId, productId uuid.UUID) error
	Clear(ctx context.Context, userId uuid.UUID) error
}
